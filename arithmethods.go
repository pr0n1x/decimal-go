package dec

import "math/big"

func (d Decimal) lhs() *DecimalMut {
	if d.p == nil {
		return &DecimalMut{exp: 0, val: big.Int{}}
	}
	return d.p.Copy()
}

func (d Decimal) Add(rhs Decimal) Decimal {
	return d.lhs().Add(rhs).Val()
}

func (d Decimal) Sub(rhs Decimal) Decimal {
	return d.lhs().Sub(rhs).Val()
}

func (d Decimal) Mul(rhs Decimal) Decimal {
	return d.lhs().Mul(rhs).Val()
}

func (d Decimal) Quo(rhs Decimal) Decimal {
	return d.lhs().Quo(rhs).Val()
}

func (d Decimal) QuoRem(rhs Decimal) (quo, rem Decimal) {
	qm, rm := d.lhs().QuoRem(rhs, nil)
	return qm.Val(), rm.Val()
}

func (d Decimal) Div(rhs Decimal) Decimal {
	return d.lhs().Div(rhs).Val()
}

func (d Decimal) Mod(rhs Decimal) Decimal {
	return d.lhs().Mod(rhs).Val()
}

func (d Decimal) DivMod(rhs Decimal) (div, mod Decimal) {
	dm, tm := d.lhs().DivMod(rhs, d.p.exp.Zero().Var())
	return dm.Val(), tm.Val()
}

func (d Decimal) QuoTail(rhs Decimal) (Decimal, Decimal) {
	div, tail := d.lhs().QuoTail(rhs, nil)
	return div.Val(), tail.Val()
}

func (d Decimal) DivTail(rhs Decimal) (Decimal, Decimal) {
	div, tail := d.lhs().DivTail(rhs, nil)
	return div.Val(), tail.Val()
}

func (d Decimal) Abs() Decimal {
	return d.lhs().Abs().Val()
}

func (d Decimal) Neg() Decimal {
	return d.lhs().Neg().Val()
}

func (d Decimal) Cmp(rhs Decimal) int {
	// reimplement coercePrecision to skip probable copying.
	if rhs.p == nil {
		rhs.p = &DecimalMut{}
	}
	lhs := d.lhs()
	if lhs.exp != rhs.p.exp {
		if lhs.exp < rhs.p.exp {
			lhs = lhs.Copy().Rescale(rhs.p.exp)
		} else {
			rhs = rhs.Copy().Rescale(lhs.exp)
		}
	}
	lhs.coercePrecision(&rhs)
	return lhs.val.Cmp(&rhs.p.val)
}

func (d Decimal) Round(r Precision, m RoundingMode) Decimal {
	return d.lhs().Round(r, m).Val()
}
