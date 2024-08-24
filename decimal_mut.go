package dec

import (
	"github.com/pr0n1x/go-type-wrappers/werr"
	"math/big"
)

type DecimalMut struct {
	exp Precision
	val big.Int
}

var ErrInvalidDecimalString = werr.New("invalid decimal value string")

// NewDecimalMut creates *DecimalMut from a raw *big.Int value and a precision
func NewDecimalMut(val *big.Int, precision Precision) *DecimalMut {
	d := DecimalMut{
		exp: precision,
		val: big.Int{},
	}
	d.val.Set(val)
	return &d
}

func (d *DecimalMut) Value() Decimal {
	return Decimal{p: d}
}

func (d *DecimalMut) Rescale(p Precision) *DecimalMut {
	if d == nil {
		return nil
	}
	if p > d.exp {
		d.val.Mul(&d.val, (p - d.exp).Multiplier())
	} else if p < d.exp {
		multiplier := (d.exp - p).Multiplier()
		if d.val.Sign() < 0 {
			mod := (&big.Int{}).Set(&d.val)
			mod.Abs(mod)
			mod.Mod(mod, multiplier)
			d.val.Add(&d.val, mod)
		} else {
			d.val.Set(&d.val)
		}
		d.val.Div(&d.val, multiplier)
	}
	d.exp = p
	return d
}

func (d *DecimalMut) Copy() *DecimalMut {
	if d == nil {
		return nil
	}
	r := DecimalMut{
		exp: d.exp,
		val: big.Int{},
	}
	r.val.Set(&d.val)
	return &r
}
