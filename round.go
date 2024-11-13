package dec

type RoundingMode uint8

const (
	HalfEven     RoundingMode = iota // == IEEE 754-2008 roundTiesToEven.
	HalfUp                           // == IEEE 754-2008 roundTiesToAway.
	HalfDown                         // no IEEE 754-2008 equivalent.
	ToZero                           // == IEEE 754-2008 roundTowardZero.
	AwayFromZero                     // no IEEE 754-2008 equivalent.
	// TODO: Implement other rounding modes

	//ToPositiveInf                     // == IEEE 754-2008 roundTowardPositive.
)

func (d *DecimalMut) Round(r Precision, m RoundingMode) *DecimalMut {
	switch m {
	case HalfEven, HalfUp, HalfDown, ToZero, AwayFromZero:
	default:
		panic("invalid rounding mode")
	}
	sign := d.Val().Sign()
	if d.exp <= r || sign == 0 {
		return d
	}
	rounding := d.Copy()
	remainder := rounding.RescaleRem(r)
	unit := r.Unit()
	switch m {
	case HalfEven, HalfUp, HalfDown:
		half := Unit(r + 1).Mul(FromUInt64(5, 0))
		if sign < 0 {
			half.Var().Neg()
		}
		halfDeflection := remainder.Cmp(half)
		switch {
		case halfDeflection == 0 && m == HalfEven:
			if rounding.Copy().Mod(FromUInt64(2, 0)).Val().Int64() != 0 {
				rounding.Add(unit)
			}
		case sign > 0 && ((halfDeflection == 0 && m == HalfUp) || halfDeflection > 0):
			rounding.Add(unit)
		case sign < 0 && ((halfDeflection == 0 && m == HalfDown) || halfDeflection < 0):
			rounding.Sub(unit)
		}
	case ToZero, AwayFromZero:
		if remainder.Var().Abs().Val().Cmp(Zero(0)) > 0 && m == AwayFromZero {
			if sign > 0 {
				rounding.Add(unit)
			} else {
				rounding.Sub(unit)
			}
		}
	}
	d.val = rounding.val
	d.exp = rounding.exp
	return d
}
