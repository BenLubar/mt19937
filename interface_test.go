package mt19937_test

import (
	"encoding"
	"math/rand"

	"github.com/BenLubar/mt19937"
)

var _ rand.Source = (*mt19937.Source)(nil)
var _ encoding.BinaryMarshaler = (*mt19937.Source)(nil)
var _ encoding.BinaryUnmarshaler = (*mt19937.Source)(nil)
