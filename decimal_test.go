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
	if got, expected := n.Int64(), int64(-2_203_000_000); got != expected {
		t.Errorf("MustParse(\"-2.203\").Int64() != %d, got %d", expected, got)
	}
}

func TestNegativeLtZero(t *testing.T) {
	n := Nano.MustParse("-1.23")
	if got, expected := n.Div(Nano.FromUInt64(10)).String(), "-0.123"; got != expected {
		t.Errorf("wrong negative value lower than zero: expected %q, got %q", expected, got)
	}
	if got, expected := n.Div(Nano.FromUInt64(100)).String(), "-0.0123"; got != expected {
		t.Errorf("wrong negative value lower than zero: expected %q, got %q", expected, got)
	}
}
