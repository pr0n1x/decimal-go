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
	var (
		precisionDelta Precision
		lessPrecise    *Decimal
		morePrecise    *Decimal
	)
	if a.precision > b.precision {
		precisionDelta = a.precision - b.precision
		morePrecise = a
		lessPrecise = b
		side = 1
	} else {
		precisionDelta = b.precision - a.precision
		morePrecise = b
		lessPrecise = a
		side = -1
	}

	precisionMultiplier := (&big.Int{}).Exp(
		big.NewInt(10),
		(&big.Int{}).SetUint64(uint64(precisionDelta)),
		nil)

	increasedPrecisionValue := big.Int{}
	increasedPrecisionValue.Mul(lessPrecise.value, precisionMultiplier)
	if checkNumberSize(&increasedPrecisionValue) {
		lessPrecise.precision += precisionDelta
		lessPrecise.value = &increasedPrecisionValue
	} else {
		morePrecise.precision -= precisionDelta
		reducedPrecisionValue := big.Int{}
		reducedPrecisionValue.Div(morePrecise.value, precisionMultiplier)
		if !checkNumberSize(&reducedPrecisionValue) {
			panic("unreachable: more precise value can't be too big after precision reduction")
		}
		morePrecise.value = &reducedPrecisionValue
		side *= -1
	}
	return side
}

func Add(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	return MustFromUnits((&big.Int{}).Add(a.value, b.value), a.precision)
}

func Sub(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	return MustFromUnits((&big.Int{}).Sub(a.value, b.value), a.precision)
}

func Mul(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	value := &big.Int{}
	value.Mul(a.value, b.value)
	value.Div(value, a.precision.Multiplier())
	return MustFromUnits(value, a.precision)
}

func Div(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	value := &big.Int{}
	value.Mul(a.value, a.precision.Multiplier())
	value.Div(value, b.value)
	return MustFromUnits(value, a.precision)
}

// DivTail returns a division result and a tail (residual/remainder related to a precision)
func DivTail(a, b Decimal) (Decimal, Decimal) {
	coercePrecision(&a, &b)
	value := &big.Int{}
	tail := &big.Int{}
	value.Mul(a.value, a.precision.Multiplier())
	value.DivMod(value, b.value, tail)
	tail.Div(tail, (a.precision - 1).Multiplier())
	return MustFromUnits(value, a.precision), MustFromUnits(tail, a.precision+1)
}

func Mod(a, b Decimal) Decimal {
	coercePrecision(&a, &b)
	return MustFromUnits((&big.Int{}).Mod(a.value, b.value), a.precision)
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
	neg.value.Neg(a.value)
	neg.precision = a.precision
	return neg
}

func Cmp(a, b Decimal) int {
	coercePrecision(&a, &b)
	return a.value.Cmp(b.value)
}
