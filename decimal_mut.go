package dec

import (
	"math/big"

	"github.com/pr0n1x/go-liners/werr"
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
		d.val.Mul(&d.val, (p - d.exp).multiplierOnlyForReadIPromise())
	} else if p < d.exp {
		multiplier := (d.exp - p).multiplierOnlyForReadIPromise()
		remainder.p = &DecimalMut{exp: d.exp, val: big.Int{}}
		rem := &remainder.p.val
		d.val.QuoRem(&d.val, multiplier, rem)
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
