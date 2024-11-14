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
	d.val.Div(&d.val, d.exp.multiplierOnlyForReadIPromise())
	return d
}

func (d *DecimalMut) Quo(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Mul(&d.val, d.exp.multiplierOnlyForReadIPromise())
	d.val.Quo(&d.val, &rhs.p.val)
	return d
}

func (d *DecimalMut) QuoRem(rhs Decimal, rem *DecimalMut) (*DecimalMut, *DecimalMut) {
	d.coercePrecision(&rhs)
	if rem == nil {
		rem = &DecimalMut{exp: d.exp, val: *big.NewInt(0)}
	} else {
		*rem = DecimalMut{exp: d.exp, val: *big.NewInt(0)}
	}
	d.val.QuoRem(&d.val, &rhs.p.val, &rem.val)
	d.val.Mul(&d.val, d.exp.multiplierOnlyForReadIPromise())
	return d, rem
}

func (d *DecimalMut) Div(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Mul(&d.val, d.exp.multiplierOnlyForReadIPromise())
	rem := big.Int{}
	d.val.QuoRem(&d.val, &rhs.p.val, &rem)
	if sign := rem.Sign(); sign != 0 {
		one := big.Int{} // on stack
		one.SetInt64(int64(sign))
		d.val.Add(&d.val, &one)
	}
	return d
}

// DivTail returns a division result and a tail (residual/remainder related to a precision).
// For the operation `res, tail := x.DivTail(y)`
// there is a valid equation `x = (res - tail) * y`.
// e.g. for operation using Milli precision:
// `res, tail := Milli.FromUint64(2).DivTail(Milli.FromUint64(3))`,
// result and tail are:
// res = 0.666
// tail = 0.002,
// 0.666 * 2 = 1.998
// 1.998 + 0.002 = 2.
func (d *DecimalMut) DivTail(rhs Decimal, tail *DecimalMut) (*DecimalMut, *DecimalMut) {
	d.coercePrecision(&rhs)
	if tail == nil {
		tail = &DecimalMut{exp: d.exp * 2, val: *big.NewInt(0)}
	} else {
		*tail = DecimalMut{exp: d.exp * 2, val: *big.NewInt(0)}
	}
	d.val.Mul(&d.val, d.exp.multiplierOnlyForReadIPromise())
	d.val.QuoRem(&d.val, &rhs.p.val, &tail.val)
	return d, tail
}

func (d *DecimalMut) Mod(rhs Decimal) *DecimalMut {
	d.coercePrecision(&rhs)
	d.val.Mod(&d.val, &rhs.p.val)
	return d
}

func (d *DecimalMut) DivMod(rhs Decimal, m *DecimalMut) (*DecimalMut, *DecimalMut) {
	d.coercePrecision(&rhs)
	precisionMultiplier := d.exp.multiplierOnlyForReadIPromise()
	var a, b big.Int
	a.Mul(&d.val, precisionMultiplier)
	b.Mul(&rhs.p.val, precisionMultiplier)
	d.val.DivMod(&a, &b, &m.val)
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
