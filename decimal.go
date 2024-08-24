package dec

import (
	"github.com/pr0n1x/go-type-wrappers/werr"
	"math/big"
	"strings"
)

var (
	ErrInvalidDecimalString = werr.New("invalid decimal value string")
)

// Decimal based on tlb.Coins from tonutils-go
type Decimal struct {
	value     *big.Int
	precision Precision
}

// TODO: add methods Ceil, Floor, Round, Pow, Avg(first Decimal, rest ...Decimal).

func Zero(p Precision) Decimal {
	return FromUnitsUInt64(0, p)
}

func One(p Precision) Decimal {
	return FromUInt64(1, p)
}

func Ten(p Precision) Decimal {
	return FromUInt64(10, p)
}

func Unit(p Precision) Decimal {
	return FromUnits((&big.Int{}).SetUint64(1), p)
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

func (d Decimal) Rescale(p Precision) Decimal {
	if p > d.precision {
		return FromUnits((&big.Int{}).Mul(d.value, (p-d.precision).Multiplier()), p)
	}
	if p < d.precision {
		multiplier := (d.precision - p).Multiplier()
		value := &big.Int{}
		if d.value.Sign() < 0 {
			mod := (&big.Int{}).Set(d.value)
			mod.Abs(mod)
			mod.Mod(mod, multiplier)
			value.Add(d.value, mod)
		} else {
			value.Set(d.value)
		}
		return FromUnits(value.Div(value, multiplier), p)
	}
	return d.Copy()
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

func (d Decimal) UInt64() uint64 {
	if d.value == nil {
		return 0
	}
	return d.value.Uint64()
}

func (d Decimal) Int64() int64 {
	if d.value == nil {
		return 0
	}
	return d.value.Int64()
}

// FromUnits creates Decimal from a raw *big.Int value and a precision
func FromUnits(val *big.Int, precision Precision) Decimal {
	return Decimal{
		precision: precision,
		value:     (&big.Int{}).Set(val),
	}
}

// FromUnitsUInt64 creates Decimal from a raw uint64 value and a precision
func FromUnitsUInt64(val uint64, precision Precision) Decimal {
	return FromUnits((&big.Int{}).SetUint64(val), precision)
}

// FromUnitsInt64 creates Decimal from a raw int64 value and a precision
func FromUnitsInt64(val int64, precision Precision) Decimal {
	return FromUnits((&big.Int{}).SetInt64(val), precision)
}

// FromUInt64 creates Decimal using uint64 as an integer part of the value
func FromUInt64(val uint64, precision Precision) Decimal {
	value := (&big.Int{}).SetUint64(val)
	value.Mul(value, precision.Multiplier())
	return Decimal{value: value, precision: precision}
}

// FromInt64 creates Decimal using int64 as an integer part of the value
func FromInt64(val int64, precision Precision) Decimal {
	value := (&big.Int{}).SetInt64(val)
	value.Mul(value, precision.Multiplier())
	return Decimal{value: value, precision: precision}
}

// Parse parses decimal number
func Parse(val string, precision Precision) (Decimal, error) {
	s := strings.SplitN(val, ".", 2)

	if len(s) == 0 {
		return Decimal{}, ErrInvalidDecimalString
	}

	hi, ok := (&big.Int{}).SetString(s[0], 10)
	if !ok {
		return Decimal{}, ErrInvalidDecimalString
	}

	hi = hi.Mul(hi, (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil))

	if len(s) == 2 {
		loStr := s[1]
		// lo can have max {decimals} digits
		if len(loStr) > int(precision) {
			loStr = loStr[:precision]
		}

		leadZeroes := 0
		for _, sym := range loStr {
			if sym != '0' {
				break
			}
			leadZeroes++
		}

		lo, ok := (&big.Int{}).SetString(loStr, 10)
		if !ok {
			return Decimal{}, ErrInvalidDecimalString
		}

		digits := len(lo.String()) // =_=
		lo = lo.Mul(lo, (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64((int(precision)-leadZeroes)-digits)), nil))
		if hi.Sign() < 0 && lo.Sign() > 0 {
			hi = hi.Sub(hi, lo)
		} else {
			hi = hi.Add(hi, lo)
		}
	}

	return FromUnits(hi, precision), nil
}

// MustParse the same as ParsePrecise but panics on error
func MustParse(val string, precision Precision) Decimal {
	return must(Parse(val, precision))
}

// ParseUnits parse a string of whole number containing integer and fractional part of the value
func ParseUnits(val string, precision Precision) (Decimal, error) {
	if bn, ok := (&big.Int{}).SetString(val, 10); ok {
		return FromUnits(bn, precision), nil
	}
	return Decimal{}, ErrInvalidDecimalString
}

// MustParseUnits the same as ParseUnits but panics on error
func MustParseUnits(val string, precision Precision) Decimal {
	return must(ParseUnits(val, precision))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
