package benchmarks

import (
	"math/big"
	"testing"

	dec "github.com/pr0n1x/decimal-go"
)

func Benchmark_Decimal_Immutable(b *testing.B) {
	million := dec.Quecto.MustParse("1000000")
	divider := dec.Quecto.MustParse("33.67")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = million.Div(divider)
	}
}

func Benchmark_Decimal_Mutable(b *testing.B) {
	millions := make([]dec.Decimal, 1_000_000)
	for i := range millions {
		millions[i] = dec.Quecto.MustParse("1000000")
	}

	divider := dec.Quecto.MustParse("33.67")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		millions[i].Var().Div(divider)
	}
}

func Benchmark_BigRat_Immutable(b *testing.B) {
	divider := &big.Rat{}
	divider.SetString("33.67")
	million := &big.Rat{}
	million.SetString("1000000")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		(&big.Rat{}).Set(million).Quo(million, divider)
	}
}

func Benchmark_BigRat_Mutable(b *testing.B) {
	million := &big.Rat{}
	million.SetString("1000000")
	millions := make([]*big.Rat, 1_000_000)
	for i := range millions {
		millions[i] = new(big.Rat).Set(million)
	}
	divider := &big.Rat{}
	divider.SetString("33.67")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		millions[i].Quo(million, divider)
	}
}
