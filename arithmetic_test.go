package dec

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
)

type testFrac struct {
	n Decimal // n
	d Decimal // d
}

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

func Test_Quo_Fractional(t *testing.T) {
	res := Nano.MustParse("10").Quo(Nano.MustParse("0.5"))
	if res, expected := res.String(), "20"; res != expected {
		t.Fatalf("invalid multiplication, expected %s, got %s", expected, res)
	}
}

func Test_Div_FractionRound(t *testing.T) {
	res := Nano.FromUInt64(2).Div(Nano.FromUInt64(3))
	if got, expected := res.String(), "0.666666667"; got != expected {
		t.Fatalf("invalid division result last precision digit round: expected %s, got %s", expected, got)
	}
}

func Test_Mod_Fractional(t *testing.T) {
	res := Nano.MustParse("0.1").Mod(Nano.MustParse("0.03"))
	if res, expected := res.String(), "0.01"; res != expected {
		t.Fatalf("invalid multiplication, expected %s, got %s", expected, res)
	}
}

func Test_DivMod(t *testing.T) {
	for _, tc := range []struct {
		frac testFrac
		res  string
		mod  string
	}{
		{frac: testFrac{n: Centi.MustParse("5"), d: Centi.MustParse("3")}, res: "1", mod: "2"},
		{frac: testFrac{n: Centi.MustParse("6.67"), d: Centi.MustParse("3.3")}, res: "2", mod: "0.07"},
		{frac: testFrac{n: Milli.FromUInt64(2), d: Milli.FromUInt64(3)}, res: "0", mod: "2"},
		{frac: testFrac{n: Nano.MustParse("464"), d: Nano.MustParse("33")}, res: "14", mod: "2"},
		{frac: testFrac{n: Nano.MustParse("4.64"), d: Nano.MustParse("3.3")}, res: "1", mod: "1.34"},
		{frac: testFrac{n: Milli.MustParse("6.6"), d: Milli.MustParse("3")}, res: "2", mod: "0.6"},
		{frac: testFrac{n: Milli.FromUInt64(5), d: Milli.FromUInt64(3)}, res: "1", mod: "2"},
	} {
		{
			div, rem := tc.frac.n.QuoRem(tc.frac.d)
			if got, expected := div.String(), tc.res; got != expected {
				t.Fatalf(`invalid "div" of QuoRem operation, expected %s, got %s`, expected, got)
			}
			if got, expected := rem.String(), tc.mod; got != expected {
				t.Fatalf(`invalid "rem" of QuoRem operation, expected %s, got %s`, expected, got)
			}
		}
		{
			div, mod := tc.frac.n.DivMod(tc.frac.d)
			if got, expected := div.String(), tc.res; got != expected {
				t.Fatalf(`invalid "div" of DivMod operation, expected %s, got %s`, expected, got)
			}
			if got, expected := mod.String(), tc.mod; got != expected {
				t.Fatalf(`invalid "mod" of DivMod operation, expected %s, got %s`, expected, got)
			}
		}
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
	for _, frac := range []testFrac{
		{n: MustParse("1.004", Milli), d: MustParse("0.6", Milli)},
		{n: MustParse("1.004", Milli), d: MustParse("0.06", Milli)},
		{n: MustParse("1.000000004", Quecto), d: MustParse("0.6", Quecto)},
		{n: MustParse("39.999999999999999999999999999999", Quecto), d: MustParse("131072", Quecto)},
		{n: MustParse("1.2", Milli), d: MustParse("0.6", Milli)},
		{n: MustParse("1.998", Milli), d: MustParse("3", Milli)},
		{n: MustParse("2", Milli), d: MustParse("3", Milli)},
	} {
		numerator, denominator := frac.n, frac.d
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
