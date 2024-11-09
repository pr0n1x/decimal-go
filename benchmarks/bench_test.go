package benchmarks

import (
	"math/big"
	"testing"

	dec "github.com/pr0n1x/decimal-go"
)

func Benchmark_Immutable(b *testing.B) {
	million := dec.Quecto.MustParse("1000000")
	divider := dec.Quecto.MustParse("33.67")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = million.Div(divider)
	}
}

func Benchmark_Mutable(b *testing.B) {
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

func Benchmark_Rat_Mmute(b *testing.B) {
	million := &big.Rat{}
	million.SetString("1000000")
	divider := &big.Rat{}
	divider.SetString("33.67")

	ress := make([]big.Rat, 1_000_000)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ress[i].Quo(million, divider)
	}
}
