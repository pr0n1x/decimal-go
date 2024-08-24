package dec

func (d Decimal) Add(rhs Decimal) Decimal {
	return d.p.Copy().Add(rhs).Value()
}

func (d Decimal) Sub(rhs Decimal) Decimal {
	return d.p.Copy().Sub(rhs).Value()
}

func (d Decimal) Mul(rhs Decimal) Decimal {
	return d.p.Copy().Mul(rhs).Value()
}

func (d Decimal) Div(rhs Decimal) Decimal {
	return d.p.Copy().Div(rhs).Value()
}

func (d Decimal) Mod(rhs Decimal) Decimal {
	return d.p.Copy().Mod(rhs).Value()
}

func (d Decimal) DivMod(rhs Decimal) (div, mod Decimal) {
	dm, tm := d.p.Copy().DivMod(rhs, d.p.exp.Zero().Mutable())
	return dm.Value(), tm.Value()
}

func (d Decimal) DivTail(rhs Decimal) (div, tail Decimal) {
	dm, tm := d.p.Copy().DivTail(rhs, d.p.exp.Zero().Mutable())
	return dm.Value(), tm.Value()
}

func (d Decimal) Abs() Decimal {
	return d.p.exp.Zero().Mutable().Abs(d).Value()
}

func (d Decimal) Neg() Decimal {
	return d.p.exp.Zero().Mutable().Neg(d).Value()
}

func (d Decimal) Cmp(rhs Decimal) int {
	// reimplement coercePrecision to skip probable copying
	if rhs.p == nil {
		rhs.p = &DecimalMut{}
	}
	lhs := d.p
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
	return d.p.Copy().Round(r, m).Value()
}
