package dec

import "math/big"

var (
	zMultiplier, _      = (&big.Int{}).SetString("1", BASE)
	deciMultiplier, _   = (&big.Int{}).SetString("10", BASE)
	centiMultiplier, _  = (&big.Int{}).SetString("100", BASE)
	milliMultiplier, _  = (&big.Int{}).SetString("1000", BASE)
	microMultiplier, _  = (&big.Int{}).SetString("1000000", BASE)
	nanoMultiplier, _   = (&big.Int{}).SetString("1000000000", BASE)
	picoMultiplier, _   = (&big.Int{}).SetString("1000000000000", BASE)
	femtoMultiplier, _  = (&big.Int{}).SetString("1000000000000000", BASE)
	attoMultiplier, _   = (&big.Int{}).SetString("1000000000000000000", BASE)
	zeptoMultiplier, _  = (&big.Int{}).SetString("1000000000000000000000", BASE)
	yoctoMultiplier, _  = (&big.Int{}).SetString("1000000000000000000000000", BASE)
	rontoMultiplier, _  = (&big.Int{}).SetString("1000000000000000000000000000", BASE)
	quectoMultiplier, _ = (&big.Int{}).SetString("1000000000000000000000000000000", BASE)
	multiplierCache     = make(map[Precision]*big.Int)
)

func (p Precision) multiplierPromiseReadOnly() *big.Int {
	switch p {
	case Z:
		return zMultiplier
	case Deci:
		return deciMultiplier
	case Centi:
		return centiMultiplier
	case Milli:
		return milliMultiplier
	case Micro:
		return microMultiplier
	case Nano:
		return nanoMultiplier
	case Pico:
		return picoMultiplier
	case Femto:
		return femtoMultiplier
	case Atto:
		return attoMultiplier
	case Zepto:
		return zeptoMultiplier
	case Yocto:
		return yoctoMultiplier
	case Ronto:
		return rontoMultiplier
	case Quecto:
		return quectoMultiplier
	}
	if value, ok := multiplierCache[p]; ok {
		return value
	}
	value := big.Int{}
	value.SetUint64(BASE)
	value.Exp(&value, big.NewInt(int64(p)), nil)
	multiplierCache[p] = &value
	return &value
}