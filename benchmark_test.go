package mt19937_test

import (
	"testing"

	"github.com/BenLubar/mt19937"
	"github.com/BenLubar/mt19937/internal/mt19937_c"
)

func BenchmarkRand(b *testing.B) {
	b.Run("Go", func(b *testing.B) {
		var s mt19937.Source
		s.Seed(123456789)
		for i := 0; i < b.N; i++ {
			s.Uint64()
		}
	})
	b.Run("C", mt19937_c.BenchmarkRand)
}

func BenchmarkSeed(b *testing.B) {
	b.Run("Go", func(b *testing.B) {
		var s mt19937.Source
		for i := 0; i < b.N; i++ {
			s.Seed(123456789)
		}
	})
	b.Run("C", mt19937_c.BenchmarkSeed)
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
	b.Run("Go", func(b *testing.B) {
		var s mt19937.Source
		for i := 0; i < b.N; i++ {
			s.SeedMulti(seed[:]...)
		}
	})
	b.Run("C", mt19937_c.BenchmarkSeedMulti)
}
