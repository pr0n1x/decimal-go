package dec

import "math/big"

type Precision uint8

// https://www.nist.gov/pml/owm/metric-si-prefixes
const (
	Milli Precision = (1 + iota) * 3
	Micro
	Nano
	Pico
	Femto
	Atto
	Zepto
	Yocto
	Ronto
	Quecto
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

func (p Precision) FromUnits(val *big.Int) (Decimal, error) {
	return FromUnits(val, p)
}

func (p Precision) MustFromUnits(val *big.Int) Decimal {
	return MustFromUnits(val, p)
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
