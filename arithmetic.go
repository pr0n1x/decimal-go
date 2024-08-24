package dec

import (
	"math/big"
)

func coercePrecision(a, b *Decimal) (side int8) {
	if a.precision == b.precision {
		return 0
	}
	if a.value == nil {
		a.value = &big.Int{}
	}
	if b.value == nil {
		b.value = &big.Int{}
	}
	if a.precision > b.precision {
		*b = b.Rescale(a.precision)
		side = 1
	} else {
		*a = a.Rescale(b.precision)
		side = -1
	}
	return side
}

func Add(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	return FromUnits((&big.Int{}).Add(a.value, b.value), a.precision)
}

func Sub(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	return FromUnits((&big.Int{}).Sub(a.value, b.value), a.precision)
}

func Mul(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	value := &big.Int{}
	value.Mul(a.value, b.value)
	value.Div(value, a.precision.Multiplier())
	return FromUnits(value, a.precision)
}

func Div(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	value := &big.Int{}
	value.Mul(a.value, a.precision.Multiplier())
	value.Div(value, b.value)
	return FromUnits(value, a.precision)
}

// DivTail returns a division result and a tail (residual/remainder related to a precision)
func DivTail(a, b Decimal) (Decimal, Decimal) {
	coercePrecision(&a, &b)
	value := &big.Int{}
	tail := &big.Int{}
	value.Mul(a.value, a.precision.Multiplier())
	value.DivMod(value, b.value, tail)
	tail.Div(tail, (a.precision - 1).Multiplier())
	return FromUnits(value, a.precision), FromUnits(tail, a.precision+1)
}

func Mod(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	return FromUnits((&big.Int{}).Mod(a.value, b.value), a.precision)
}

func DivMod(a, b Decimal) (div, mod Decimal) {
	coercePrecision(&a, &b)
	div.value = &big.Int{}
	mod.value = &big.Int{}
	div.precision = a.precision
	mod.precision = a.precision
	aValue := &big.Int{}
	bValue := &big.Int{}
	precisionMultiplier := a.precision.Multiplier()
	aValue.Mul(a.value, precisionMultiplier)
	bValue.Mul(b.value, precisionMultiplier)
	div.value, mod.value = (&big.Int{}).DivMod(aValue, bValue, mod.value)
	div.value.Mul(div.value, precisionMultiplier)
	mod.value.Div(mod.value, precisionMultiplier)
	return div, mod
}

func Abs(signed Decimal) (absolute Decimal) {
	absolute.value = (&big.Int{}).Abs(signed.value)
	absolute.precision = signed.precision
	return absolute
}

func Neg(a Decimal) (neg Decimal) {
	neg.value = &big.Int{}
	neg.value.Neg(a.value)
	neg.precision = a.precision
	return neg
}

func Cmp(a, b Decimal) int {
	coercePrecision(&a, &b)
	return a.value.Cmp(b.value)
}

type RoundingMode uint8

const (
	HalfEven RoundingMode = iota // == IEEE 754-2008 roundTiesToEven
	HalfUp                       // == IEEE 754-2008 roundTiesToAway
	HalfDown                     // no IEEE 754-2008 equivalent
	// TODO: Implement other rounding modes
	//ToZero                            // == IEEE 754-2008 roundTowardZero
	//AwayFromZero                      // no IEEE 754-2008 equivalent
	//ToNegativeInf                     // == IEEE 754-2008 roundTowardNegative
	//ToPositiveInf                     // == IEEE 754-2008 roundTowardPositive
)

func Round(d Decimal, r Precision, m RoundingMode) Decimal {
	switch m {
	case HalfEven, HalfUp, HalfDown:
	default:
		panic("invalid rounding mode")
	}
	sign := d.Units().Sign()
	if d.precision <= r || sign == 0 {
		return d.Copy()
	}
	rounded := d.Rescale(r)
	unit := r.Unit()
	switch m {
	case HalfEven, HalfUp, HalfDown:
		half := Unit(r + 1).Mul(FromUInt64(5, 0))
		mod := d.Abs().Rescale(r + 1).Mod(unit)
		halfDeflection := mod.Cmp(half)
		if halfDeflection == 0 && m == HalfEven {
			if rounded.Mod(FromUInt64(2, 0)).Int64() != 0 {
				rounded = rounded.Add(unit)
			}
		} else if (halfDeflection == 0 && m == HalfUp) || halfDeflection > 0 {
			if sign > 0 {
				rounded = rounded.Add(unit)
			} else {
				rounded = rounded.Sub(unit)
			}
		}
	}
	return rounded
}
