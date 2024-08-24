package dec

import "testing"

func Test_RoundHalfUp(t *testing.T) {
	n := Nano.MustParse("1.123456789")
	if got, expected := n.Round(5, HalfUp).String(), "1.12346"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
	if got, expected := n.Round(4, HalfUp).String(), "1.1235"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
	if got, expected := n.Round(3, HalfUp).String(), "1.123"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
	if got, expected := n.Round(0, HalfUp).String(), "1"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
}

func Test_RoundHalfUpNeg(t *testing.T) {
	n := Nano.MustParse("-1.123456789")
	if got, expected := n.Round(5, HalfUp).String(), "-1.12346"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
	if got, expected := n.Round(4, HalfUp).String(), "-1.1235"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
	if got, expected := n.Round(3, HalfUp).String(), "-1.123"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
	if got, expected := n.Round(0, HalfUp).String(), "-1"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
}

func Test_RoundHalfDown(t *testing.T) {
	n := Nano.MustParse("1.123456789")
	if got, expected := n.Round(4, HalfDown).String(), "1.1234"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
}

func Test_RoundHalfDownNeg(t *testing.T) {
	n := Nano.MustParse("-1.123456789")
	if got, expected := n.Round(4, HalfDown).String(), "-1.1234"; got != expected {
		t.Fatalf("invalid Round, expected %s, got %s", expected, got)
	}
}
