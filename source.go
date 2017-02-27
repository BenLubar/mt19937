// The original copyright notice is reproduced below:

/*
   A C-program for MT19937-64 (2004/9/29 version).
   Coded by Takuji Nishimura and Makoto Matsumoto.

   This is a 64-bit version of Mersenne Twister pseudorandom number
   generator.

   Before using, initialize the state by using init_genrand64(seed)
   or init_by_array64(init_key, key_length).

   Copyright (C) 2004, Makoto Matsumoto and Takuji Nishimura,
   All rights reserved.

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:

     1. Redistributions of source code must retain the above copyright
        notice, this list of conditions and the following disclaimer.

     2. Redistributions in binary form must reproduce the above copyright
        notice, this list of conditions and the following disclaimer in the
        documentation and/or other materials provided with the distribution.

     3. The names of its contributors may not be used to endorse or promote
        products derived from this software without specific prior written
        permission.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR
   CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
   EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
   PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
   PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
   LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
   NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
   SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

   References:
   T. Nishimura, ``Tables of 64-bit Mersenne Twisters''
     ACM Transactions on Modeling and
     Computer Simulation 10. (2000) 348--357.
   M. Matsumoto and T. Nishimura,
     ``Mersenne Twister: a 623-dimensionally equidistributed
       uniform pseudorandom number generator''
     ACM Transactions on Modeling and
     Computer Simulation 8. (Jan. 1998) 3--30.

   Any feedback is very welcome.
   http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/emt.html
   email: m-mat @ math.sci.hiroshima-u.ac.jp (remove spaces)
*/

// Package mt19937 is a pure Go port of the C reference implementation of a
// 64-bit Mersenne Twister random number generator, originally written by
// Takuji Nishimura and Makoto Matsumoto.
//
// http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/VERSIONS/C-LANG/mt19937-64.c
package mt19937

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

const nn = 312
const mm = 156
const matrixA = 0xb5026f5aa96619e9
const um = 0xffffffff80000000
const lm = 0x7fffffff

// Source is a rand.Source implementation using a 64 bit Mersenne Twister
// pseudorandom number generator.
//
// Source is not safe for use from multiple concurrent goroutines.
type Source struct {
	state [nn]uint64
	index uint16
	init  bool
}

var _ rand.Source64 = (*Source)(nil)

// NewSource is a convenience function that creates a Source and calls Seed.
func NewSource(seed int64) *Source {
	var s Source
	s.Seed(seed)
	return &s
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *Source) MarshalBinary() ([]byte, error) {
	if !s.init {
		return nil, nil
	}

	var buf [nn*64/8 + 16/8]byte
	for i, n := range s.state {
		binary.LittleEndian.PutUint64(buf[i*64/8:], n)
	}
	binary.LittleEndian.PutUint16(buf[nn*64/8:], s.index)
	return buf[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Source) UnmarshalBinary(buf []byte) error {
	if len(buf) == 0 {
		s.init = false
		return nil
	}

	if len(buf) != nn*64/8+16/8 {
		return fmt.Errorf("mt19937: random state size mismatch: %d != %d", nn*64/8+16/8, len(buf))
	}

	index := binary.LittleEndian.Uint16(buf[nn*64/8:])
	if index >= nn {
		return fmt.Errorf("mt19937: random state index out of range: %d", index)
	}

	s.init = true
	for i := range s.state {
		s.state[i] = binary.LittleEndian.Uint64(buf[i*64/8:])
	}
	s.index = index
	return nil
}

// Seed resets the random seed to a specific value.
func (s *Source) Seed(seed int64) {
	s.init = true
	s.state[0] = uint64(seed)
	for s.index = 1; s.index < nn; s.index++ {
		s.state[s.index] = (6364136223846793005*(s.state[s.index-1]^(s.state[s.index-1]>>62)) + uint64(s.index))
	}
}

// SeedMulti resets the seed to a larger value if more entropy is desired.
func (s *Source) SeedMulti(seed ...uint64) {
	var i, j, k uint64
	s.Seed(19650218)
	i = 1
	j = 0
	if nn > len(seed) {
		k = nn
	} else {
		k = uint64(len(seed))
	}
	for ; k != 0; k-- {
		s.state[i] = (s.state[i] ^ ((s.state[i-1] ^ (s.state[i-1] >> 62)) * 3935559000370003845)) + seed[j] + j
		i++
		j++
		if i >= nn {
			s.state[0] = s.state[nn-1]
			i = 1
		}
		if j >= uint64(len(seed)) {
			j = 0
		}
	}
	for k = nn - 1; k != 0; k-- {
		s.state[i] = (s.state[i] ^ ((s.state[i-1] ^ (s.state[i-1] >> 62)) * 2862933555777941757)) - i
		i++
		if i >= nn {
			s.state[0] = s.state[nn-1]
			i = 1
		}
	}

	s.state[0] = 1 << 63
}

// Int63 returns a uniformly random integer in the range [0, 2^63).
func (s *Source) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Uint64 returns a uniformly random integer in the range [0, 2^64).
func (s *Source) Uint64() uint64 {
	if !s.init {
		s.Seed(5489)
	}

	if s.index >= nn {
		mag01 := [2]uint64{0, matrixA}

		var i int
		for i = 0; i < nn-mm; i++ {
			x := (s.state[i] & um) | (s.state[i+1] & lm)
			s.state[i] = s.state[i+mm] ^ (x >> 1) ^ mag01[x&1]
		}
		for ; i < nn-1; i++ {
			x := (s.state[i] & um) | (s.state[i+1] & lm)
			s.state[i] = s.state[i+(mm-nn)] ^ (x >> 1) ^ mag01[x&1]
		}
		x := (s.state[nn-1] & um) | (s.state[0] & lm)
		s.state[nn-1] = s.state[mm-1] ^ (x >> 1) ^ mag01[x&1]

		s.index = 0
	}

	x := s.state[s.index]
	s.index++

	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)

	return x
}
