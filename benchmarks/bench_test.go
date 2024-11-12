package benchmarks

import (
	"math/big"
	"testing"

	dec "github.com/pr0n1x/decimal-go"
	"github.com/pr0n1x/go-liners/assert"
)

const testPrecision = dec.Nano

func Benchmark_Immutable_Decimal(b *testing.B) {
	million := testPrecision.MustParse("1000000")
	divider := testPrecision.MustParse("33.67")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = million.Div(divider)
	}
}

func Benchmark_Immutable_BigRat(b *testing.B) {
	divider := &big.Rat{}
	divider.SetString("33.67")
	million := &big.Rat{}
	million.SetString("1000000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		(&big.Rat{}).Set(million).Quo(million, divider)
	}
}

func Benchmark_Mutable_Decimal(b *testing.B) {
	millions := make([]dec.Decimal, 1_000_000)
	for i := range millions {
		millions[i] = testPrecision.MustParse("1000000")
	}

	divider := testPrecision.MustParse("33.67")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		millions[i].Var().Div(divider)
	}
}

func Benchmark_Mutable_BigRat(b *testing.B) {
	million := &big.Rat{}
	million.SetString("1000000")
	millions := make([]*big.Rat, 1_000_000)
	for i := range millions {
		millions[i] = new(big.Rat).Set(million)
	}
	divider := &big.Rat{}
	divider.SetString("33.67")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		millions[i].Quo(million, divider)
	}
}

const testBigFrac = "9183278754354899275983457827.98723867459876125012387684587367745647"
const testBigFracPrecision = 38

func Benchmark_BigFraction_Div_Decimal(b *testing.B) {
	rat := dec.MustParse(testBigFrac, testBigFracPrecision)
	three := dec.FromUInt64(3, testBigFracPrecision)
	results := make([]dec.Decimal, 1_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = append(results, rat.Div(three))
	}
}

func Benchmark_BigFraction_Div_BigRat(b *testing.B) {
	rat := assert.Ok(new(big.Rat).SetString(testBigFrac))
	three := new(big.Rat).SetInt64(3)
	results := make([]*big.Rat, 1_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = append(results, new(big.Rat).Quo(rat, three))
	}
}

func Benchmark_BigFraction_Mul_Decimal(b *testing.B) {
	rat := dec.MustParse(testBigFrac, testBigFracPrecision)
	three := dec.FromUInt64(3, testBigFracPrecision)
	results := make([]dec.Decimal, 1_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = append(results, rat.Mul(three))
	}
}

func Benchmark_BigFraction_Mul_BigRat(b *testing.B) {
	rat := assert.Ok(new(big.Rat).SetString(testBigFrac))
	three := new(big.Rat).SetInt64(3)
	results := make([]*big.Rat, 1_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = append(results, new(big.Rat).Mul(rat, three))
	}
}
