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
