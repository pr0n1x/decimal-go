package dec

import "testing"

type roundTestCase struct {
	n Decimal   // number.
	r Precision // round precision.
	e string    // expected.
}

func Test_RoundHalfUp(t *testing.T) {
	const (
		ns  = "1.123456789"
		nns = "-1.123456789"
	)
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse(ns), r: 5, e: "1.12346"},
		{n: Nano.MustParse(ns), r: 4, e: "1.1235"},
		{n: Nano.MustParse(ns), r: 3, e: "1.123"},
		{n: Nano.MustParse(ns), r: 0, e: "1"},
		{n: Nano.MustParse(nns), r: 5, e: "-1.12346"},
		{n: Nano.MustParse(nns), r: 4, e: "-1.1235"},
		{n: Nano.MustParse(nns), r: 3, e: "-1.123"},
		{n: Nano.MustParse(nns), r: 0, e: "-1"},
	} {
		if got, expected := tc.n.Round(tc.r, HalfUp).String(), tc.e; got != expected {
			t.Fatalf("invalid Round, expected %s, got %s", expected, got)
		}
	}
}

func Test_RoundHalfDown(t *testing.T) {
	for _, tc := range []roundTestCase{
		{n: Nano.MustParse("1.123456789"), r: 4, e: "1.1234"},
		{n: Nano.MustParse("-1.123456789"), r: 4, e: "-1.1234"},
	} {
		if got, expected := tc.n.Round(tc.r, HalfDown).String(), tc.e; got != expected {
			t.Fatalf("invalid Round, expected %s, got %s", expected, got)
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
		if got, expected := tc.n.Round(tc.r, AwayFromZero).String(), tc.e; got != expected {
			t.Fatalf("invalid Round, expected %s, got %s", expected, got)
		}
	}
}
