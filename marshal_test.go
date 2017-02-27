package mt19937_test

import (
	"testing"
	"testing/quick"

	"github.com/BenLubar/mt19937"
)

func TestMarshal(t *testing.T) {
	if err := quick.CheckEqual(func(seed int64, skip uint8) (out [count]uint64) {
		s := mt19937.NewSource(seed)
		for i := 0; i < int(skip); i++ {
			s.Uint64()
		}
		for i := range out {
			out[i] = s.Uint64()
		}
		return
	}, func(seed int64, skip uint8) (out [count]uint64) {
		s := mt19937.NewSource(seed)
		for i := 0; i < int(skip); i++ {
			s.Uint64()
		}
		b, err := s.MarshalBinary()
		if err != nil {
			t.Errorf("for seed %d: %v", seed, err)
		}
		s = new(mt19937.Source)
		err = s.UnmarshalBinary(b)
		if err != nil {
			t.Errorf("for seed %d: %v", seed, err)
		}
		for i := range out {
			out[i] = s.Uint64()
		}
		return
	}, nil); err != nil {
		t.Error(err)
	}
}

func TestMarshalEmpty(t *testing.T) {
	var a mt19937.Source
	if b, err := a.MarshalBinary(); err != nil {
		t.Error(err)
	} else if len(b) != 0 {
		t.Errorf("expected empty: % x", b)
	}
	var b mt19937.Source
	b.Seed(123456789)
	if err := b.UnmarshalBinary(nil); err != nil {
		t.Error(err)
	}

	for i := 0; i < count; i++ {
		if j, k := a.Uint64(), b.Uint64(); j != k {
			t.Errorf("%d: %d != %d", i, j, k)
		}
	}
}

func TestUnmarshalError(t *testing.T) {
	t.Run("Size", func(t *testing.T) {
		b, err := mt19937.NewSource(12345).MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		var s mt19937.Source
		if err = s.UnmarshalBinary(b); err != nil {
			t.Error(err)
		}
		if err = s.UnmarshalBinary(b[:len(b)-1]); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("Index", func(t *testing.T) {
		b, err := mt19937.NewSource(12345).MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		var s mt19937.Source
		if err = s.UnmarshalBinary(b); err != nil {
			t.Error(err)
		}
		b[len(b)-1] = 255
		if err = s.UnmarshalBinary(b); err == nil {
			t.Error("expected error")
		}
	})
}
