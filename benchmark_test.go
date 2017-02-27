package mt19937_test

import (
	"testing"

	"github.com/BenLubar/mt19937"
	"github.com/BenLubar/mt19937/internal/mt19937_c"
)

func init() {
	// so that C benchmarks don't reduce test coverage
	mt19937_c.BenchmarkRand(&testing.B{N: 1})
	mt19937_c.BenchmarkSeed(&testing.B{N: 1})
	mt19937_c.BenchmarkSeedMulti(&testing.B{N: 1})
}

func BenchmarkRand(b *testing.B) {
	b.Run("C", mt19937_c.BenchmarkRand)
	b.Run("Go", func(b *testing.B) {
		var s mt19937.Source
		s.Seed(123456789)
		for i := 0; i < b.N; i++ {
			s.Uint64()
		}
	})
	b.Run("CGo", func(b *testing.B) {
		mt19937_c.InitGenRand64(123456789)
		for i := 0; i < b.N; i++ {
			mt19937_c.GenRand64Int64()
		}
	})
}

func BenchmarkSeed(b *testing.B) {
	b.Run("C", mt19937_c.BenchmarkSeed)
	b.Run("Go", func(b *testing.B) {
		var s mt19937.Source
		for i := 0; i < b.N; i++ {
			s.Seed(123456789)
		}
	})
	b.Run("CGo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mt19937_c.InitGenRand64(123456789)
		}
	})
}

func BenchmarkSeedMulti(b *testing.B) {
	seed := [8]uint64{
		0x1234567890123456,
		0x7890123456789012,
		0x3456789012345678,
		0x9012345678901234,
		0x5678901234567890,
		0x1234567890123456,
		0x7890123456789012,
		0x3456789012345678,
	}
	b.Run("C", mt19937_c.BenchmarkSeedMulti)
	b.Run("Go", func(b *testing.B) {
		var s mt19937.Source
		for i := 0; i < b.N; i++ {
			s.SeedMulti(seed[:]...)
		}
	})
	b.Run("CGo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mt19937_c.InitByArray64(seed[:]...)
		}
	})
}
