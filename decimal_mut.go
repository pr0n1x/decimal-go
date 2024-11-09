package dec

import (
	"math/big"

	"github.com/pr0n1x/go-type-wrappers/werr"
)

type DecimalMut struct {
	exp Precision
	val big.Int
}

var ErrInvalidDecimalString = werr.New("invalid decimal value string")

// NewDecimalMut creates *DecimalMut from a raw *big.Int value and a precision.
func NewDecimalMut(val *big.Int, precision Precision) *DecimalMut {
	d := DecimalMut{
		exp: precision,
		val: big.Int{},
	}
	d.val.Set(val)
	return &d
}

func (d *DecimalMut) Val() Decimal {
	return Decimal{p: d}
}

func (d *DecimalMut) RescaleRem(p Precision) (remainder Decimal) {
	if d == nil {
		return
	}
	if p > d.exp {
		d.val.Mul(&d.val, (p - d.exp).Multiplier())
	} else if p < d.exp {
		multiplier := (d.exp - p).Multiplier()
		remainder.p = &DecimalMut{exp: d.exp, val: big.Int{}}
		rem := &remainder.p.val
		if d.val.Sign() < 0 {
			rem.Set(&d.val)
			rem.Abs(rem)
			rem.Mod(rem, multiplier)
			d.val.Add(&d.val, rem)
			d.val.Div(&d.val, multiplier)
		} else {
			d.val.DivMod(&d.val, multiplier, rem)
		}
	}
	d.exp = p
	return
}

func (d *DecimalMut) Rescale(p Precision) *DecimalMut {
	d.RescaleRem(p)
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
