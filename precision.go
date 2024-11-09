package dec

import (
	"math/big"
)

type Precision uint16

// https://www.nist.gov/pml/owm/metric-si-prefixes
const (
	Z      Precision = 0
	Deci   Precision = 1
	Centi  Precision = 2
	Milli  Precision = 3
	Micro  Precision = 6
	Nano   Precision = 9
	Pico   Precision = 12
	Femto  Precision = 15
	Atto   Precision = 18
	Zepto  Precision = 21
	Yocto  Precision = 24
	Ronto  Precision = 27
	Quecto Precision = 30
)

func (p Precision) Increase(delta Precision) (Precision, bool) {
	const maxUint16 = 1<<16 - 1 // 65535.
	if delta > maxUint16 || maxUint16-delta < p {
		return p, false
	}
	return p + delta, true
}

func (p Precision) Decrease(delta Precision) (Precision, bool) {
	if delta > p {
		return p, false
	}
	return p - delta, true
}

func (p Precision) Zero() Decimal { return Zero(p) }

func (p Precision) One() Decimal { return One(p) }

func (p Precision) Ten() Decimal { return Ten(p) }

func (p Precision) Multiplier() *big.Int { return PrecisionMultiplier(p) }

func (p Precision) Unit() Decimal { return Unit(p) }

func (p Precision) MaxFraction() Decimal { return One(p + 1).Sub(Unit(p + 1)).Rescale(p) }

func (p Precision) FromUnits(val *big.Int) Decimal { return FromUnits(val, p) }

func (p Precision) FromUnitsUInt64(val uint64) Decimal { return FromUnitsUInt64(val, p) }

func (p Precision) FromUnitsInt64(val int64) Decimal { return FromUnitsInt64(val, p) }

func (p Precision) FromUInt64(val uint64) Decimal { return FromUInt64(val, p) }

func (p Precision) FromInt64(val int64) Decimal { return FromInt64(val, p) }

func (p Precision) Parse(val string) (Decimal, error) { return Parse(val, p) }

func (p Precision) MustParse(val string) Decimal { return MustParse(val, p) }

func (p Precision) ParseUnits(val string) (Decimal, error) { return ParseUnits(val, p) }

func (p Precision) MustParseUnits(val string) Decimal { return MustParseUnits(val, p) }
