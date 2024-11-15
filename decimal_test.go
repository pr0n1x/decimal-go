package dec

import (
	"math/big"
	"testing"
)

func Test_Parse(t *testing.T) {
	if Nano.MustParse("1.1").Units().Uint64() != 1_100_000_000 {
		t.Fatal(`MustParse("1.1").Units().Uint64() != 1_100_000_000`)
	}
}

func Test_ParseNano(t *testing.T) {
	if Nano.MustParseUnits("1100000000").Units().Uint64() != 1_100_000_000 {
		t.Fatal(`MustParseNano("1100000000").Units().Uint64() != 1_100_000_000`)
	}
}

func Test_FromBigInt(t *testing.T) {
	var n int64 = 2_200_999_888
	if Nano.FromUnits(big.NewInt(n)).Units().Int64() != n {
		t.Fatal("FromUnits(big.NewInt(n)).Units().Int64() != n")
	}
}

func TestParseNeg(t *testing.T) {
	n, err := Nano.Parse("-2.203")
	if err != nil {
		t.Fatal(err)
	}
	if got, expected := n.Units().Int64(), int64(-2_203_000_000); got != expected {
		t.Errorf("MustParse(\"-2.203\").Int64() != %d, got %d", expected, got)
	}
}

func TestNegativeLtZero(t *testing.T) {
	n := Nano.MustParse("-1.23")
	if got, expected := n.Quo(Nano.FromUInt64(10)).String(), "-0.123"; got != expected {
		t.Errorf("wrong negative value lower than zero: expected %q, got %q", expected, got)
	}
	if got, expected := n.Quo(Nano.FromUInt64(100)).String(), "-0.0123"; got != expected {
		t.Errorf("wrong negative value lower than zero: expected %q, got %q", expected, got)
	}
}

func TestToUInt64(t *testing.T) {
	n := uint64(123)
	d := Nano.FromUInt64(n)
	if d.UInt64() != n {
		t.Errorf("wrong u64 value; expected %d, got %d", n, d.UInt64())
	}
}

func TestToInt64(t *testing.T) {
	n := int64(123)
	d := Nano.FromInt64(n)
	if d.Int64() != n {
		t.Errorf("wrong u64 value; expected %d, got %d", n, d.UInt64())
	}
}

func TestDecimalMutNilPointer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("operation on nil *DecimalMut pointer should panic")
		}
	}()
	var nilPtr *DecimalMut = nil
	nilPtr.Add(Z.FromUInt64(2))
}

func TestRescaleReminder(t *testing.T) {
	for _, tc := range []struct {
		number    Decimal
		rescaleTo Precision
		rescaled  string
		remainder string
	}{
		{number: Milli.MustParse("123.456"), rescaleTo: Deci, rescaled: "1234", remainder: "0.056"},
		{number: Milli.MustParse("-123.456"), rescaleTo: Deci, rescaled: "-1234", remainder: "-0.056"},
		{number: Milli.MustParse("-123.456"), rescaleTo: Micro, rescaled: "-123456000", remainder: "0"},
	} {
		rescaled, remainder := tc.number.RescaleRem(tc.rescaleTo)
		if got, expected := rescaled.Units().String(), tc.rescaled; got != expected {
			t.Fatalf("invalid rescaled value: expected '%s', got '%s'", expected, got)
		}
		if got, expected := remainder.Precision(), tc.number.Precision(); got != expected {
			t.Fatal("invalid remainder precision")
		}
		if got, expected := remainder.String(), tc.remainder; got != expected {
			t.Fatalf("invalid rescale reminder: expected '%s', got '%s'", expected, got)
		}
	}
}

func TestParseLimitedPrecision(t *testing.T) {
	for _, tc := range []struct {
		n  string    // number string
		p  Precision // possible precision, depends from autoPrecision
		li bool      // limited precision
		ep Precision // expected precision
		eu string    // expected units value
	}{
		{n: "123.456", p: 0, li: false, ep: Milli, eu: "123456"},
		{n: "123.456", p: Centi, li: false, ep: Milli, eu: "123456"},
		{n: "123.456", p: Centi, li: true, ep: Centi, eu: "12345"},
		{n: "123.456", p: Deci, li: true, ep: Deci, eu: "1234"},
		{n: "123.456", p: Micro, li: false, ep: Micro, eu: "123456000"},
		{n: "123.456789", p: 0, li: false, ep: Micro, eu: "123456789"},
		{n: "123.456789", p: Milli, li: true, ep: Milli, eu: "123456"},
	} {
		n := must(Parse(tc.n, tc.p, tc.li))
		if got, expected := n.Units().String(), tc.eu; got != expected {
			t.Fatalf("invalid parsed units valus: expected '%s', got '%s'", expected, got)
		}
	}
}
