package dec

import (
	"math/big"
	"testing"
)

func Test_MaxFitValues(t *testing.T) {
	for _, size := range []FitSize{Fit128, Fit256} {
		maxValue := Nano.FromUnits(size.MaxValue())
		if _, fit := maxValue.Fit(size); !fit {
			t.Fatalf("max value does not fit into %d bytes", size)
		}
		maxPlusOne := (&big.Int{}).Add(maxValue.Units(), big.NewInt(1))
		if _, fit := Nano.FromUnits(maxPlusOne).Fit(size); fit {
			t.Fatalf("max value + 1 should not fit into %d bytes", size)
		}
	}
}

func Test_NonFitOperations(t *testing.T) {
	for _, size := range []FitSize{Fit128, Fit256} {
		maxValue := Nano.FromUnits(size.MaxValue())
		maxFit := maxValue.MustFit(size)
		unitFit := Nano.Unit().MustFit(size)
		if _, fit := maxFit.Add(unitFit); fit {
			t.Fatalf("max value + 1 should not fit into %d bytes", size)
		}
		twoFit := Nano.FromUInt64(2).MustFit(size)
		if _, fit := maxFit.Div(twoFit); !fit {
			t.Fatalf("max value / 2 should fit into %d bytes", size)
		}
	}
}
