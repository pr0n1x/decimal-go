package dec

import (
	"github.com/pr0n1x/go-type-wrappers/werr"
	"math/big"
	"strings"
)

var (
	ErrInvalid          = werr.New("invalid string")
	ErrInvalidPrecision = werr.New("invalid precision")
	ErrTooBigNumber     = werr.New("too big number")
)

// Decimal based on tlb.Coins from tonutils-go
type Decimal struct {
	value     *big.Int
	precision Precision
}

func Zero(precision Precision) Decimal {
	return MustFromUInt64(0, precision)
}

func One(precision Precision) Decimal {
	return MustFromZUInt64(1, precision)
}

func Ten(precision Precision) Decimal {
	return MustFromZUInt64(10, precision)
}

// Precision precision - max decimals digits
func (d Decimal) Precision() Precision {
	return d.precision
}

// Units raw big int value
func (d Decimal) Units() *big.Int {
	if d.value == nil {
		return big.NewInt(0)
	}
	return (&big.Int{}).Set(d.value)
}

func (d Decimal) Copy() Decimal {
	return Decimal{
		precision: d.precision,
		value:     d.Units(),
	}
}

func (d Decimal) String() string {
	if d.value == nil {
		return "0"
	}

	a := d.value.String()
	if a == "0" {
		// process 0 faster and simpler
		return a
	}

	splitter := len(a) - int(d.precision)
	if splitter <= 0 {
		a = "0." + strings.Repeat("0", int(d.precision)-len(a)) + a
	} else {
		// set . between lo and hi
		a = a[:splitter] + "." + a[splitter:]
	}

	// cut last zeroes
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == '.' {
			a = a[:i]
			break
		}
		if a[i] != '0' {
			a = a[:i+1]
			break
		}
	}

	return a
}
