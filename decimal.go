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

// TODO: exclude checkNumberSize from creation and rescaling
// TODO: implement method .Fit(bits uint8) (Decimal, error) which rescales number (with the same logic as in checkNumberSize)
// TODO: change IncreasePrecision to Rescale method which can increase and reduce precision
// TODO: re-implement coercePrecision using Rescale method
// TODO: add methods Ceil, Floor, Round, Pow, Avg(first Decimal, rest ...Decimal).

func Zero(p Precision) Decimal {
	return MustFromUInt64(0, p)
}

func One(p Precision) Decimal {
	return MustFromZUInt64(1, p)
}

func Ten(p Precision) Decimal {
	return MustFromZUInt64(10, p)
}

func Unit(p Precision) Decimal {
	return MustFromUnits((&big.Int{}).SetUint64(1), p)
}

func PrecisionMultiplier(p Precision) *big.Int {
	var value big.Int
	value.SetUint64(10)
	value.Exp(&value, big.NewInt(int64(p)), nil)
	return &value
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

// IncreasePrecision increases value precision
func (d Decimal) IncreasePrecision(p Precision) (Decimal, error) {
	precision := d.precision + 1
	value := (&big.Int{}).Mul(d.value, p.Multiplier())
	return FromUnits(value, precision)
}

// MustIncreasePrecision the same as IncreasePrecision but panics on error
func (d Decimal) MustIncreasePrecision(p Precision) Decimal {
	return must(d.IncreasePrecision(p))
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
