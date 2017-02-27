// +build go1.8

package mt19937_test

import (
	"math/rand"

	"github.com/BenLubar/mt19937"
)

var _ rand.Source64 = (*mt19937.Source)(nil)
