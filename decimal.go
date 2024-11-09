package dec

import (
	"math/big"
	"strings"
)

// Decimal based on tlb.Coins from tonutils-go.
type Decimal struct {
	p *DecimalMut
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

func (d Decimal) Var() *DecimalMut {
	return d.p
}

// Precision exp - max decimals digits.
func (d Decimal) Precision() Precision {
	return d.p.exp
}

// Units raw big int.
func (d Decimal) Units() *big.Int {
	if d.p == nil {
		return big.NewInt(0)
	}
	return (&big.Int{}).Set(&d.p.val)
}

func (d Decimal) Sign() int {
	if d.p == nil {
		return 0
	}
	return d.p.val.Sign()
}

func (d Decimal) Rescale(p Precision) Decimal {
	return d.p.Copy().Rescale(p).Val()
}

func (d Decimal) RescaleRem(p Precision) (rescaled, remainder Decimal) {
	rescaled = d.p.Copy().Val()
	remainder = rescaled.p.RescaleRem(p)
	return rescaled, remainder
}

func (d Decimal) Copy() Decimal {
	return d.p.Copy().Val()
}

func (d Decimal) String() string {
	if d.p == nil {
		return "0"
	}

	sign := d.p.val.Sign()
	if sign == 0 {
		// process 0 faster and simpler.
		return "0"
	}
	a := d.p.val.String()
	if sign < 0 {
		a = a[1:]
	}
	splitter := len(a) - int(d.p.exp)
	if splitter <= 0 {
		a = "0." + strings.Repeat("0", int(d.p.exp)-len(a)) + a
	} else {
		// set . between lo and hi.
		a = a[:splitter] + "." + a[splitter:]
	}

	// cut last zeroes.
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

	if sign < 0 {
		a = "-" + a
	}
	return a
}

func (d Decimal) UInt64() uint64 {
	if d.p == nil {
		return 0
	}
	return d.p.val.Div(&d.p.val, d.p.exp.Multiplier()).Uint64()
}

func (d Decimal) Int64() int64 {
	if d.p == nil {
		return 0
	}
	return d.p.val.Div(&d.p.val, d.p.exp.Multiplier()).Int64()
}

// FromUnits creates Decimal from a raw *big.Int value and a precision.
func FromUnits(val *big.Int, precision Precision) Decimal {
	return Decimal{p: NewDecimalMut(val, precision)}
}

// FromUnitsUInt64 creates Decimal from a raw uint64 value and a precision.
func FromUnitsUInt64(val uint64, precision Precision) Decimal {
	return FromUnits((&big.Int{}).SetUint64(val), precision)
}

// FromUnitsInt64 creates Decimal from a raw int64 value and a precision.
func FromUnitsInt64(val int64, precision Precision) Decimal {
	return FromUnits((&big.Int{}).SetInt64(val), precision)
}

// FromUInt64 creates Decimal using uint64 as an integer part of the value.
func FromUInt64(val uint64, precision Precision) Decimal {
	value := (&big.Int{}).SetUint64(val)
	value.Mul(value, precision.Multiplier())
	return Decimal{p: NewDecimalMut(value, precision)}
}

// FromInt64 creates Decimal using int64 as an integer part of the value.
func FromInt64(val int64, precision Precision) Decimal {
	value := (&big.Int{}).SetInt64(val)
	value.Mul(value, precision.Multiplier())
	return Decimal{p: NewDecimalMut(value, precision)}
}

// Parse parses decimal number.
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
		// lo can have max {decimals} digits.
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

		digits := len(lo.String()) // =_=.
		lo = lo.Mul(lo, (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64((int(precision)-leadZeroes)-digits)), nil))
		if hi.Sign() < 0 && lo.Sign() > 0 {
			hi = hi.Sub(hi, lo)
		} else {
			hi = hi.Add(hi, lo)
		}
	}

	return FromUnits(hi, precision), nil
}

// MustParse the same as ParsePrecise but panics on error.
func MustParse(val string, precision Precision) Decimal {
	return must(Parse(val, precision))
}

// ParseUnits parse a string of whole number containing integer and fractional part of the value.
func ParseUnits(val string, precision Precision) (Decimal, error) {
	if bn, ok := (&big.Int{}).SetString(val, 10); ok {
		return FromUnits(bn, precision), nil
	}
	return Decimal{}, ErrInvalidDecimalString
}

// MustParseUnits the same as ParseUnits but panics on error.
func MustParseUnits(val string, precision Precision) Decimal {
	return must(ParseUnits(val, precision))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
