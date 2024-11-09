package dec

import (
	"testing"

	"github.com/pr0n1x/go-type-wrappers/assert"
)

func Test_RescaleUp(t *testing.T) {
	for _, testCase := range []struct {
		num Decimal
		bef string
		aft string
	}{
		{num: Deci.MustParse("123.4"), bef: "1234", aft: "1234000"},
		{num: Deci.MustParse("654.3"), bef: "6543", aft: "6543000"},
		{num: Deci.MustParse("-123.4"), bef: "-1234", aft: "-1234000"},
		{num: Deci.MustParse("-654.3"), bef: "-6543", aft: "-6543000"},
	} {
		n, unitsBefore, unitsAfter := testCase.num, testCase.bef, testCase.aft
		if n.Units().String() != unitsBefore {
			t.Fatal("invalid units")
		}
		up := n.Copy()
		delta := Precision(3)
		remainder := up.Var().RescaleRem(assert.Ok(n.Precision().Increase(delta)))
		if got := up.Units().String(); got != unitsAfter {
			t.Fatalf("invalid rescale up: expected units after rescale %s, got %s", unitsAfter, got)
		}
		if remainder.Cmp(Zero(0)) != 0 {
			t.Fatal("invalid rescale up: remainder should be zero")
		}
	}
}

func Test_RescaleDown(t *testing.T) {
	for _, testCase := range []struct {
		num Decimal
		bef string
		aft string
		rem string
	}{
		{num: Milli.MustParse("123.456"), bef: "123456", aft: "1234", rem: "0.056"},
		{num: Milli.MustParse("654.321"), bef: "654321", aft: "6543", rem: "0.021"},
		{num: Milli.MustParse("-123.456"), bef: "-123456", aft: "-1234", rem: "0.056"},
		{num: Milli.MustParse("-654.321"), bef: "-654321", aft: "-6543", rem: "0.021"},
	} {
		n, unitsBefore, unitsAfter, expectedRemainder := testCase.num, testCase.bef, testCase.aft, testCase.rem
		if n.Units().String() != unitsBefore {
			t.Fatal("invalid units")
		}
		down := n.Copy()
		delta := Precision(2)
		remainder := down.Var().RescaleRem(assert.Ok(n.Precision().Decrease(delta)))
		if got := down.Units().String(); got != unitsAfter {
			t.Fatalf("invalid rescale down: expected units after rescale %s, got %s", unitsAfter, got)
		}
		if got := remainder.String(); got != expectedRemainder {
			t.Fatalf("invalid rescale down: expected remainder %s, got %s", expectedRemainder, got)
		}
	}
}
