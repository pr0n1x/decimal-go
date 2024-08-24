package dec

import (
	"math/big"
)

func (d *DecimalMut) coercePrecision(rhs *Decimal) *DecimalMut {
	if rhs.p == nil {
		rhs.p = &DecimalMut{}
	}
	if d.exp == rhs.p.exp {
		return d
	}
	if d.exp > rhs.p.exp {
		*rhs = rhs.Copy().Rescale(d.exp)
	} else {
		*d = *d.Copy().Rescale(rhs.p.exp)
	}
	return d
}

func (d *DecimalMut) Set(val Decimal) *DecimalMut {
	if d == nil {
		return NewDecimalMut(val.Units(), val.Precision())
	}
	d.exp = val.p.exp
	d.val = big.Int{}
	d.val.Set(&val.p.val)
	return d
}

func (d *DecimalMut) Add(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Add(&d.val, &rhs.p.val)
	return d
}

func (d *DecimalMut) Sub(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Sub(&d.val, &rhs.p.val)
	return d
}

func (d *DecimalMut) Mul(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Mul(&d.val, &rhs.p.val)
	d.val.Div(&d.val, d.exp.Multiplier())
	return d
}

func (d *DecimalMut) Div(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Mul(&d.val, d.exp.Multiplier())
	d.val.Div(&d.val, &rhs.p.val)
	return d
}

// DivTail returns a division result and a tail (residual/remainder related to a precision)
func (d *DecimalMut) DivTail(rhs Decimal, tail *DecimalMut) (*DecimalMut, *DecimalMut) {
	d.coercePrecision(&rhs)
	if tail == nil {
		tail = &DecimalMut{exp: 0, val: *big.NewInt(0)}
	}
	d.val.Mul(&d.val, d.exp.Multiplier())
	d.val.DivMod(&d.val, &rhs.p.val, &tail.val)
	tail.val.Div(&tail.val, (d.exp - 1).Multiplier())
	tail.exp += 1
	return d, tail
}

func (d *DecimalMut) Mod(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Mod(&d.val, &rhs.p.val)
	return d
}

func (d *DecimalMut) DivMod(rhs Decimal, m *DecimalMut) (*DecimalMut, *DecimalMut) {
	d.coercePrecision(&rhs)
	precisionMultiplier := d.exp.Multiplier()
	aValue := (&big.Int{}).Mul(&d.val, precisionMultiplier)
	bValue := (&big.Int{}).Mul(&rhs.p.val, precisionMultiplier)
	d.val.DivMod(aValue, bValue, &m.val)
	d.val.Mul(&d.val, precisionMultiplier)
	m.val.Div(&m.val, precisionMultiplier)
	return d, m
}

func (d *DecimalMut) Abs(a Decimal) *DecimalMut {
	d.exp = a.p.exp
	d.val.Abs(&a.p.val)
	return d
}

func (d *DecimalMut) Neg(a Decimal) *DecimalMut {
	d.val.Neg(&a.p.val)
	return d
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

func (d *DecimalMut) Round(r Precision, m RoundingMode) *DecimalMut {
	switch m {
	case HalfEven, HalfUp, HalfDown:
	default:
		panic("invalid rounding mode")
	}
	sign := d.Value().Sign()
	if d.exp <= r || sign == 0 {
		return d
	}
	rounded := d.Copy().Rescale(r)
	unit := r.Unit()
	switch m {
	case HalfEven, HalfUp, HalfDown:
		half := Unit(r + 1).Mul(FromUInt64(5, 0))
		mod := d.Abs(d.Value()).Rescale(r + 1)
		mod.Mod(unit)
		halfDeflection := mod.Value().Cmp(half)
		if halfDeflection == 0 && m == HalfEven {
			if rounded.Copy().Mod(FromUInt64(2, 0)).Value().Int64() != 0 {
				rounded.Add(unit)
			}
		} else if (halfDeflection == 0 && m == HalfUp) || halfDeflection > 0 {
			if sign > 0 {
				rounded.Add(unit)
			} else {
				rounded.Sub(unit)
			}
		}
	}
	d.val = rounded.val
	d.exp = rounded.exp
	return d
}
