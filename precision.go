package dec

import "math/big"

type Precision uint8

// https://www.nist.gov/pml/owm/metric-si-prefixes
const (
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

func (p Precision) Zero() Decimal {
	return Zero(p)
}

func (p Precision) One() Decimal {
	return One(p)
}

func (p Precision) Ten() Decimal {
	return Ten(p)
}

func (p Precision) Multiplier() *big.Int {
	return PrecisionMultiplier(p)
}

func (p Precision) Unit() Decimal {
	return Unit(p)
}

func (p Precision) MaxFraction() Decimal {
	return One(p + 1).Sub(Unit(p + 1)).MustRescale(p)
}

func (p Precision) FromUnits(val *big.Int) (Decimal, error) {
	return FromUnits(val, p)
}

func (p Precision) MustFromUnits(val *big.Int) Decimal {
	return MustFromUnits(val, p)
}

func (p Precision) FromUInt64Units(val uint64) (Decimal, error) {
	return FromUInt64Units(val, p)
}

func (p Precision) MustFromUInt64Units(val uint64) Decimal {
	return MustFromUInt64Units(val, p)
}

func (p Precision) FromInt64Units(val int64) (Decimal, error) {
	return FromInt64Units(val, p)
}

func (p Precision) MustFromInt64Units(val int64) Decimal {
	return MustFromInt64Units(val, p)
}

func (p Precision) FromUInt64(val uint64) (Decimal, error) {
	return FromUInt64(val, p)
}

func (p Precision) MustFromUInt64(val uint64) Decimal {
	return MustFromUInt64(val, p)
}

func (p Precision) FromInt64(val int64) (Decimal, error) {
	return FromInt64(val, p)
}

func (p Precision) MustFromInt64(val int64) Decimal {
	return MustFromInt64(val, p)
}

func (p Precision) Parse(val string) (Decimal, error) {
	return Parse(val, p)
}

func (p Precision) MustParse(val string) Decimal {
	return MustParse(val, p)
}

func (p Precision) ParseUnits(val string) (Decimal, error) {
	return ParseUnits(val, p)
}

func (p Precision) MustParseUnits(val string) Decimal {
	return MustParseUnits(val, p)
}
