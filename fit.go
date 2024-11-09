package dec

import (
	"fmt"
	"math/big"
)

type Fit struct {
	Decimal
	Size FitSize
}
type FitSize uint32

const (
	Fit32  FitSize = 4
	Fit64  FitSize = 8
	Fit128 FitSize = 16
	Fit256 FitSize = 32
	Fit512 FitSize = 64
)

func (s FitSize) BitsLen() uint64 {
	return uint64(s) * 8
}

func (s FitSize) MaxValue() (maxValue *big.Int) {
	maxValue = (&big.Int{}).Exp(big.NewInt(2), big.NewInt(int64(s*8)), nil)
	maxValue = maxValue.Sub(maxValue, big.NewInt(1))
	return maxValue
}

type fitReduceFlag bool

const FitReduce fitReduceFlag = true

var (
	Max32BitsValue  = Fit32.MaxValue()
	Max64BitsValue  = Fit64.MaxValue()
	Max128BitsValue = Fit128.MaxValue()
	Max256BitsValue = Fit256.MaxValue()
	Max512BitsValue = Fit512.MaxValue()
)

func (d Decimal) Fit(size FitSize, reduce ...fitReduceFlag) (val Fit, fit bool) {
	val, fit = Fit{Decimal: d, Size: size}, checkNumberSize(&d.p.val, uint(size))
	// TODO: call reducePrecisionToFit.
	if !fit && len(reduce) > 0 && bool(reduce[0]) {
		return val, fit
	}
	return val, fit
}

func (d Decimal) MustFit(size FitSize, reduce ...fitReduceFlag) Fit {
	f, ok := d.Fit(size, reduce...)
	if !ok {
		panic(fmt.Sprintf("Decimal does not fit into %d bits", size.BitsLen()))
	}
	return f
}

func (f Fit) Add(rhs Fit) (Fit, bool) {
	return f.Decimal.Add(rhs.Decimal).Fit(f.Size)
}

func (f Fit) Sub(rhs Fit) (Fit, bool) {
	return f.Decimal.Sub(rhs.Decimal).Fit(f.Size)
}

func (f Fit) Mul(rhs Fit) (Fit, bool) {
	return f.Decimal.Mul(rhs.Decimal).Fit(f.Size)
}

func (f Fit) Div(rhs Fit) (Fit, bool) {
	return f.Decimal.Div(rhs.Decimal).Fit(f.Size)
}

func (f Fit) Mod(rhs Fit) (Fit, bool) {
	return f.Decimal.Mod(rhs.Decimal).Fit(f.Size)
}

func (f Fit) DivMod(rhs Fit) (div Fit, mod Fit, fit bool) {
	dd, md := f.Decimal.DivMod(rhs.Decimal)
	return fitPair(dd, md, f.Size)
}

func (f Fit) DivTail(rhs Fit) (Fit, Fit, bool) {
	dd, md := f.Decimal.DivTail(rhs.Decimal)
	return fitPair(dd, md, f.Size)
}

func (f Fit) Abs() Fit {
	return mustFit(f.Decimal.Abs().Fit(f.Size))
}

func (f Fit) Neg() (Fit, bool) {
	return f.Decimal.Neg().Fit(f.Size)
}

func (f Fit) Cmp(rhs Fit) int {
	return f.Decimal.Cmp(rhs.Decimal)
}

func (f Fit) MustAdd(rhs Fit) Fit {
	return mustFit(f.Add(rhs))
}

func (f Fit) MustSub(rhs Fit) Fit {
	return mustFit(f.Sub(rhs))
}

func (f Fit) MustMul(rhs Fit) Fit {
	return mustFit(f.Mul(rhs))
}

func (f Fit) MustDiv(rhs Fit) Fit {
	return mustFit(f.Div(rhs))
}

func (f Fit) MustDivMod(rhs Fit) (Fit, Fit) {
	return pairMustFit(f.DivMod(rhs))
}

func (f Fit) MustDivTail(rhs Fit) (Fit, Fit) {
	return pairMustFit(f.DivTail(rhs))
}

func checkNumberSize(val *big.Int, size uint) bool {
	bitsLen := val.BitLen()
	if val.Sign() < 0 {
		bitsLen += 1
	}
	if uint((bitsLen+7)>>3) > size {
		// (bitlen + 7)>>3 - the same as ceil(bitlen/8).
		return false
	}
	return true
}

const msgOpResNotFit = "operation result does not fit into 256 bits"

func mustFit(a Fit, ok bool) Fit {
	if !ok {
		panic(msgOpResNotFit)
	}
	return a
}

func pairMustFit(a Fit, b Fit, ok bool) (Fit, Fit) {
	if !ok {
		panic(msgOpResNotFit)
	}
	return a, b
}

func fitPair(a, b Decimal, size FitSize) (ar Fit, br Fit, fit bool) {
	var af, bf bool
	ar, af = a.Fit(size)
	br, bf = b.Fit(size)
	fit = af && bf
	return ar, br, fit
}

// TODO: implement reducePrecisionToFit
//func reducePrecisionToFit256().
