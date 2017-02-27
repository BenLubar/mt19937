package mt19937_test

import (
	"testing"
	"testing/quick"

	"github.com/BenLubar/mt19937"
	"github.com/BenLubar/mt19937/internal/mt19937_c"
)

const count = 128

func TestSingle64(t *testing.T) {
	if err := quick.CheckEqual(func(seed uint64) (out [count]uint64) {
		mt19937_c.InitGenRand64(seed)
		for i := range out {
			out[i] = mt19937_c.GenRand64Int64()
		}
		return
	}, func(seed uint64) (out [count]uint64) {
		var s mt19937.Source
		s.Seed(int64(seed))
		for i := range out {
			out[i] = s.Uint64()
		}
		return
	}, nil); err != nil {
		t.Error(err)
	}
}

func TestSingle63(t *testing.T) {
	if err := quick.CheckEqual(func(seed uint64) (out [count]int64) {
		mt19937_c.InitGenRand64(seed)
		for i := range out {
			out[i] = mt19937_c.GenRand64Int63()
		}
		return
	}, func(seed uint64) (out [count]int64) {
		var s mt19937.Source
		s.Seed(int64(seed))
		for i := range out {
			out[i] = s.Int63()
		}
		return
	}, nil); err != nil {
		t.Error(err)
	}
}

func TestArray64(t *testing.T) {
	if err := quick.CheckEqual(func(seed [8]uint64) (out [count]uint64) {
		mt19937_c.InitByArray64(seed[:]...)
		for i := range out {
			out[i] = mt19937_c.GenRand64Int64()
		}
		return
	}, func(seed [8]uint64) (out [count]uint64) {
		var s mt19937.Source
		s.SeedMulti(seed[:]...)
		for i := range out {
			out[i] = s.Uint64()
		}
		return
	}, nil); err != nil {
		t.Error(err)
	}
}
