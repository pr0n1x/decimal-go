package dec

import "testing"

type roundTestCase struct {
	n Decimal   // number.
	r Precision // round rescaleTo.
	e string    // expected.
}

func Test_RoundHalf(t *testing.T) {
	const (
		num = "1.123456789"
		neg = "-1.123456789"
	)
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse(num), r: 5, e: "1.12346"},
		{n: Nano.MustParse(num), r: 4, e: "1.1235"},
		{n: Nano.MustParse(num), r: 3, e: "1.123"},
		{n: Nano.MustParse(num), r: 0, e: "1"},

		{n: Nano.MustParse(neg), r: 5, e: "-1.12346"},
		{n: Nano.MustParse(neg), r: 4, e: "-1.1235"},
		{n: Nano.MustParse(neg), r: 3, e: "-1.123"},
		{n: Nano.MustParse(neg), r: 0, e: "-1"},
	} {
		if got, expected := tc.n.Round(tc.r, HalfUp).String(), tc.e; got != expected {
			t.Fatalf("invalid Round Up, expected %s, got %s", expected, got)
		}
	}
}

func Test_RoundHalfUp(t *testing.T) {
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse("1.555555555"), r: 8, e: "1.55555556"},
		{n: Nano.MustParse("1.555555555"), r: 7, e: "1.5555556"},

		{n: Nano.MustParse("-1.555555555"), r: 8, e: "-1.55555555"},
		{n: Nano.MustParse("-1.555555555"), r: 7, e: "-1.5555556"},

		{n: Milli.MustParse("1.5"), r: 0, e: "2"},
		{n: Milli.MustParse("-1.5"), r: 0, e: "-1"},
	} {
		if b := tc.e == "1"; b {
			_ = b
		}
		if got, expected := tc.n.Round(tc.r, HalfUp).String(), tc.e; got != expected {
			t.Fatalf("invalid Round Up, expected %s, got %s", expected, got)
		}
	}
}

func Test_RoundHalfDown(t *testing.T) {
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse("1.555555555"), r: 8, e: "1.55555555"},
		{n: Nano.MustParse("1.555555555"), r: 7, e: "1.5555556"},

		{n: Nano.MustParse("-1.555555555"), r: 8, e: "-1.55555556"},
		{n: Nano.MustParse("-1.555555555"), r: 7, e: "-1.5555556"},

		{n: Milli.MustParse("1.5"), r: 0, e: "1"},
		{n: Milli.MustParse("-1.5"), r: 0, e: "-2"},
	} {
		if b := tc.e == "1"; b {
			_ = b
		}
		if got, expected := tc.n.Round(tc.r, HalfDown).String(), tc.e; got != expected {
			t.Fatalf("invalid Round Up, expected %s, got %s", expected, got)
		}
	}
}

func Test_RoundModeToZero(t *testing.T) {
	const (
		ns  = "1.123456789"
		nns = "-1.123456789"
	)
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse(ns), r: 5, e: "1.12345"},
		{n: Nano.MustParse(ns), r: 4, e: "1.1234"},
		{n: Nano.MustParse(ns), r: 3, e: "1.123"},
		{n: Nano.MustParse(ns), r: 0, e: "1"},
		{n: Nano.MustParse(nns), r: 5, e: "-1.12345"},
		{n: Nano.MustParse(nns), r: 4, e: "-1.1234"},
		{n: Nano.MustParse(nns), r: 3, e: "-1.123"},
		{n: Nano.MustParse(nns), r: 0, e: "-1"},
	} {
		if got, expected := tc.n.Round(tc.r, ToZero).String(), tc.e; got != expected {
			t.Fatalf("invalid Round, expected %s, got %s", expected, got)
		}
	}
}

func Test_RoundModeAwayFromZero(t *testing.T) {
	const (
		ns  = "1.123456789"
		nns = "-1.123456789"
	)
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse(ns), r: 5, e: "1.12346"},
		{n: Nano.MustParse(ns), r: 4, e: "1.1235"},
		{n: Nano.MustParse(ns), r: 3, e: "1.124"},
		{n: Nano.MustParse(ns), r: 2, e: "1.13"},
		{n: Nano.MustParse(ns), r: 0, e: "2"},
		{n: Nano.MustParse(nns), r: 5, e: "-1.12346"},
		{n: Nano.MustParse(nns), r: 4, e: "-1.1235"},
		{n: Nano.MustParse(nns), r: 3, e: "-1.124"},
		{n: Nano.MustParse(nns), r: 2, e: "-1.13"},
		{n: Nano.MustParse(nns), r: 0, e: "-2"},
	} {
		if b := tc.e == "-1.12346"; b {
			_ = b
		}
		if got, expected := tc.n.Round(tc.r, AwayFromZero).String(), tc.e; got != expected {
			t.Fatalf("invalid Round, expected %s, got %s", expected, got)
		}
	}
}
