package dec

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func Test_Add_DifferentPrecision(t *testing.T) {
	a := MustParse("1.100001001", 9)
	b := MustParse("2.200002", 6)
	res := a.Add(b)
	if res.Units().Uint64() != 3_300_003_001 {
		t.Fatal("invalid add two numbers with different precisions")
	}
	if res.Precision() != 9 {
		t.Fatal("result should have exp of 9 decimal places")
	}
}

func Test_AddNeg(t *testing.T) {
	a := MustParse("3.300003001", 9)
	b := MustParse("-2.200002", 6)
	res := a.Add(b)
	if res.Units().Uint64() != 1_100_001_001 {
		t.Fatal("invalid add negative number")
	}
	if res.Precision() != 9 {
		t.Fatal("result should have exp of 9 decimal places")
	}
}

func Test_Mul_Fractional(t *testing.T) {
	res := Nano.MustParse("0.5").Mul(Nano.MustParse("0.5"))
	if strVal := res.String(); strVal != "0.25" {
		t.Fatal("invalid multiplication")
	}
}

func Test_Div_Fractional(t *testing.T) {
	res := Nano.MustParse("10").Div(Nano.MustParse("0.5"))
	if res, expected := res.String(), "20"; res != expected {
		t.Fatalf("invalid multiplication, expected %s, got %s", expected, res)
	}
}

func Test_Mod_Fractional(t *testing.T) {
	res := Nano.MustParse("0.1").Mod(Nano.MustParse("0.03"))
	if res, expected := res.String(), "0.01"; res != expected {
		t.Fatalf("invalid multiplication, expected %s, got %s", expected, res)
	}
}

func Test_DivMod_Fractional(t *testing.T) {
	div, mod := Nano.MustParse("464").DivMod(Nano.MustParse("33"))
	if res, expected := div.String(), "14"; res != expected {
		t.Fatalf(`invalid "div" of DivMod operation, expected %s, got %s`, expected, res)
	}
	if res, expected := mod.String(), "2"; res != expected {
		t.Fatalf(`invalid "mod" of DivMod operation, expected %s, got %s`, expected, res)
	}
}

func Test_SumHasMaxPrecision(t *testing.T) {
	sum := Milli.Zero()
	list := []Decimal{
		Deci.One(),
		Centi.One(),
		Milli.One(),
		Micro.One(),
		Nano.One(),
		Pico.One(),
		Femto.One(),
		Atto.One(),
		Zepto.One(),
		Yocto.One(),
		Ronto.One(),
		Quecto.One(),
	}
	for _, d := range list {
		sum = sum.Add(d)
	}
	if got, expected := sum.Precision(), Quecto; got != expected {
		t.Fatalf("invalid sum exp, expected %d, got %d", expected, got)
	}
	expectedSum := len(list)
	exactString, expectedString := sum.String(), strconv.Itoa(expectedSum)
	valid := false
	if exactString == expectedString {
		bigString := expectedString + strings.Repeat("0", int(Quecto))
		if bigValue, ok := (&big.Int{}).SetString(bigString, 10); ok {
			valid = sum.Units().Cmp(bigValue) == 0
		}
	}
	if !valid {
		t.Fatal("invalid sum")
	}
}

func Test_DivTail(t *testing.T) {
	type fraction struct {
		numerator   Decimal
		denominator Decimal
	}
	for _, frac := range []fraction{
		{
			numerator:   MustParse("1.004", Milli),
			denominator: MustParse("0.6", Milli),
		},
		{
			numerator:   MustParse("1.004", Milli),
			denominator: MustParse("0.06", Milli),
		},
		{
			numerator:   MustParse("1.000000004", Quecto),
			denominator: MustParse("0.6", Quecto),
		},
		{
			numerator:   MustParse("39.999999999999999999999999999999", Quecto),
			denominator: MustParse("131072", Quecto),
		},
		{
			numerator:   MustParse("1.2", Milli),
			denominator: MustParse("0.6", Milli),
		},
	} {
		numerator, denominator := frac.numerator, frac.denominator
		res, tail := numerator.DivTail(denominator)
		if 2*res.Precision() != tail.Precision() {
			t.Fatal("invalid DivTail: the tail precision should be twice the result precision")
		}
		res.Var().Rescale(res.Precision() * 2)
		if rev := res.Mul(denominator).Add(tail); rev.Cmp(numerator) != 0 {
			t.Fatalf("invalid DivTail: reversed = result * denuminator, but: %s != %s * %s",
				rev.String(), res.String(), denominator.String())
		}
	}
}

func Test_MaxFraction(t *testing.T) {
	maxFraction := Nano.MaxFraction()
	if res, expected := maxFraction.Units().String(), "999999999"; res != expected {
		t.Fatalf("invalid MaxFraction, expected %s, got %s", expected, res)
	}
}
