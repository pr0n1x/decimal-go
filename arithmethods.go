package dec

func (d Decimal) Add(rhs Decimal) Decimal {
	return Add(d, rhs)
}

func (d Decimal) Sub(rhs Decimal) Decimal {
	return Sub(d, rhs)
}

func (d Decimal) Mul(rhs Decimal) Decimal {
	return Mul(d, rhs)
}

func (d Decimal) Div(rhs Decimal) Decimal {
	return Div(d, rhs)
}

func (d Decimal) Mod(rhs Decimal) Decimal {
	return Mod(d, rhs)
}

func (d Decimal) DivMod(rhs Decimal) (div, mod Decimal) {
	return DivMod(d, rhs)
}

func (d Decimal) DivTail(rhs Decimal) (div, tail Decimal) {
	return DivTail(d, rhs)
}

func (d Decimal) Abs() Decimal {
	return Abs(d)
}

func (d Decimal) Neg() Decimal {
	return Neg(d)
}

func (d Decimal) Cmp(rhs Decimal) int {
	return Cmp(d, rhs)
}

func (d Decimal) Round(r Precision, m RoundingMode) Decimal {
	return Round(d, r, m)
}
