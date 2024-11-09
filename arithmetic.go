package dec

import (
	"math/big"
)

func (d *DecimalMut) coercePrecision(rhs *Decimal) *DecimalMut {
	if d == nil {
		panic("operation on nil *DecimalMut pointer")
	}
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

// DivTail returns a division result and a tail (residual/remainder related to a precision).
func (d *DecimalMut) DivTail(rhs Decimal, tail *DecimalMut) (*DecimalMut, *DecimalMut) {
	d.coercePrecision(&rhs)
	tailAlloc := DecimalMut{exp: d.exp * 2, val: *big.NewInt(0)}
	if tail == nil {
		tail = &tailAlloc
	} else {
		*tail = tailAlloc
	}
	d.val.Mul(&d.val, d.exp.Multiplier())
	d.val.DivMod(&d.val, &rhs.p.val, &tail.val)
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

func (d *DecimalMut) Abs() *DecimalMut {
	d.val.Abs(&d.val)
	return d
}

func (d *DecimalMut) Neg() *DecimalMut {
	d.val.Neg(&d.val)
	return d
}
