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
	rounded := d.Copy()
	remainder := rounded.RescaleRem(r)
	unit := r.Unit()
	switch m {
	case HalfEven, HalfUp, HalfDown:
		half := Unit(r + 1).Mul(FromUInt64(5, 0))
		halfDeflection := remainder.Rescale(r + 1).Cmp(half)
		if halfDeflection == 0 && m == HalfEven {
			if rounded.Copy().Mod(FromUInt64(2, 0)).Val().Int64() != 0 {
				rounded.Add(unit)
			}
		} else if (halfDeflection == 0 && m == HalfUp) || halfDeflection > 0 {
			if sign > 0 {
				rounded.Add(unit)
			} else {
				rounded.Sub(unit)
			}
		}
	case ToZero, AwayFromZero:
		if sign < 0 {
			func() {}()
		}
		if remainder.Cmp(Zero(0)) > 0 && m == AwayFromZero {
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
