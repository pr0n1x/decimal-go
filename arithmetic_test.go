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
	res := Add(a, b)
	if res.Units().Uint64() != 3_300_003_001 {
		t.Fatal("invalid add two numbers with different precisions")
	}
	if res.Precision() != 9 {
		t.Fatal("result should have precision of 9 decimal places")
	}
}

func Test_AddNeg(t *testing.T) {
	a := MustParse("3.300003001", 9)
	b := MustParse("-2.200002", 6)
	res := Add(a, b)
	if res.Units().Uint64() != 1_100_001_001 {
		t.Fatal("invalid add negative number")
	}
	if res.Precision() != 9 {
		t.Fatal("result should have precision of 9 decimal places")
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
	if sum.precision != Quecto {
		t.Fatal("invalid sum precision")
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

func Test_Rescale(t *testing.T) {
	n := Centi.MustParse("99.99")
	if n.Units().String() != "9999" {
		t.Fatal("invalid units")
	}
	if n.Rescale(n.precision+2).Units().String() != "999900" {
		t.Fatal("invalid rescale")
	}
}

func Test_DivTail(t *testing.T) {
	type fraction struct {
		numerator   Decimal
		denominator Decimal
	}
	for _, frac := range []fraction{
		{numerator: MustParse("1.004", Milli),
			denominator: MustParse("0.6", Milli)},
		{numerator: MustParse("1.000000004", Quecto),
			denominator: MustParse("0.6", Quecto)},
		{numerator: MustParse("39.999999999999999999999999999999", Quecto),
			denominator: MustParse("131072", Quecto)},
	} {
		numerator, denominator := frac.numerator, frac.denominator
		res, tail := numerator.DivTail(denominator)
		if tail.precision <= res.precision {
			t.Fatal("invalid DivTail: tail precision should be greater than result precision")
		}
		if tail.precision-res.precision > 1 {
			t.Fatal("invalid DivTail: tail and result precision difference should be equal to 1")
		}
		res = res.Rescale(res.precision + 1)
		if res.Mul(denominator).Add(tail).Cmp(numerator) != 0 {
			t.Fatal("invalid DivTail")
		}
	}
}

func Test_MaxFraction(t *testing.T) {
	maxFraction := Nano.MaxFraction()
	if res, expected := maxFraction.Units().String(), "999999999"; res != expected {
		t.Fatalf("invalid MaxFraction, expected %s, got %s", expected, res)
	}
}
