# Decimal numbers with a fixed point

A small fixed-point decimal library built on top of `math/big`.

```sh
go get github.com/pr0n1x/decimal-go
```

```go
import dec "github.com/pr0n1x/decimal-go"
```

## Features

- A number is stored as a pair: raw units (`big.Int`) and a precision (the number of decimal places after the point)
- Numbers with different precisions can operate with each other — the result is coerced to the maximum precision of the operands
- The library splits the arithmetic API into an immutable one (`Decimal`) and a mutable one (`*DecimalMut`)
- Parsing can either truncate the fractional part to the requested precision or extend the precision to fit the string
- Several division flavors: truncating (`Quo`), rounding (`Div`), with remainder (`QuoRem`/`DivMod`) and with an exact tail (`QuoTail`/`DivTail`)
- Rounding with `HalfEven`, `HalfUp`, `HalfDown`, `ToZero` and `AwayFromZero` modes
- The `Fit` type allows to check the bit size of a value or of an operation result (32/64/128/256/512 bits)
- Ready-made text marshaling types (`TextNano`, `TextMilli`, ...) for JSON and other text-based encodings

## Precision

`Precision` is the number of decimal places after the point. Constants are named
after the [SI prefixes](https://www.nist.gov/pml/owm/metric-si-prefixes):

| Constant | Decimal places |
|----------|----------------|
| `Z`      | 0              |
| `Deci`   | 1              |
| `Centi`  | 2              |
| `Milli`  | 3              |
| `Micro`  | 6              |
| `Nano`   | 9              |
| `Pico`   | 12             |
| `Femto`  | 15             |
| `Atto`   | 18             |
| `Zepto`  | 21             |
| `Yocto`  | 24             |
| `Ronto`  | 27             |
| `Quecto` | 30             |

Any other `Precision` value (it is a `uint16`) works as well.

## Use cases

### Create 1 (one) with a precision of 9 decimal places

```go
// from raw units (a big.Int of 1_000_000_000 units means 1.0 at Nano precision)
one := dec.FromUnits(new(big.Int).SetUint64(1_000_000_000), dec.Nano)
one := dec.Nano.FromUnits(new(big.Int).SetUint64(1_000_000_000))

// from raw uint64/int64 units
one := dec.FromUnitsUInt64(1_000_000_000, 9)
one := dec.FromUnitsUInt64(1_000_000_000, dec.Nano)
one := dec.FromUnitsInt64(1_000_000_000, dec.Nano)
one := dec.Nano.FromUnitsInt64(1_000_000_000)

// from a whole number
one := dec.FromUInt64(1, dec.Nano)
one := dec.FromInt64(1, dec.Nano)
one := dec.Nano.FromInt64(1)

// quite often zero precision is also useful
ten := dec.FromUInt64(10, dec.Z)
ten := dec.Z.FromUInt64(10)
// for example, to initialize a counter when the precision of the operation is not known yet
// (the sum gets the precision of the added numbers automatically)
zero := dec.Z.Zero()
zero := dec.Zero(dec.Z)
zero := dec.Z.FromUInt64(0)
```

### Parse a string

```go
n, err := dec.Nano.Parse("-2.203")     // -2.203, precision 9
n := dec.Nano.MustParse("1.1")         // panics on error, 1_100_000_000 units

// units string: an integer containing both the whole and the fractional parts
n := dec.Nano.MustParseUnits("1100000000") // 1.1

// the third argument of dec.Parse decides what to do when the string
// has more fractional digits than the requested precision:
n, err := dec.Parse("123.456789", dec.Milli, true)  // 123.456  — the tail is truncated
n, err := dec.Parse("123.456789", dec.Milli, false) // 123.456789 — precision is extended to Micro
n, err := dec.Parse("123.456", 0, false)            // 123.456  — precision is detected from the string
// Precision.Parse and Precision.MustParse always truncate (limitPrecision = true)
```

### Format as a string

`String()` prints the number without trailing zeroes in the fractional part:

```go
dec.Nano.MustParse("0.5").Mul(dec.Nano.MustParse("0.5")).String() // "0.25"
dec.Nano.FromUInt64(20).String()                                  // "20"
```

### Immutable arithmetic API

Each operation returns a new `Decimal` and leaves the operands untouched:

```go
number := dec.Z.FromInt64(123)
hundred := dec.Nano.FromUInt64(100)
three := dec.Z.FromInt64(3)
twenty := number.Mod(hundred).Sub(three) // 20
```

### Mutable arithmetic API

`Var()` returns a `*DecimalMut` view that shares the underlying value,
so operations modify the number in place:

```go
hundred := dec.Nano.FromUInt64(100)
three := dec.Deci.FromInt64(3)
number := dec.Milli.FromInt64(123) // number == 123
number.Var().Mod(hundred)          // number == 23
number.Var().Sub(three)            // number == 20
// the precision is coerced to the maximum of all operated numbers, i.e. Nano (9)
```

### Mixing precisions

The result of an operation always has the maximum precision of its operands:

```go
a := dec.MustParse("1.100001001", 9, true)
b := dec.MustParse("2.200002", 6, true)
sum := a.Add(b)   // 3.300003001
sum.Precision()   // 9
```

### Division flavors

```go
two, three := dec.Nano.FromUInt64(2), dec.Nano.FromUInt64(3)

// Quo truncates the result towards zero
two.Quo(three) // 0.666666666

// Div rounds the last digit away from zero when there is a non-zero remainder
two.Div(three) // 0.666666667

// QuoRem returns the integer quotient and the remainder
quo, rem := dec.Centi.MustParse("6.67").QuoRem(dec.Centi.MustParse("3.3")) // 2 and 0.07

// DivMod is the same, but based on Euclidean division ((big.Int).DivMod)
div, mod := dec.Centi.MustParse("5").DivMod(dec.Centi.MustParse("3")) // 1 and 2
```

`QuoTail`/`DivTail` return the division result together with an exact "tail" —
the residue that did not fit into the result precision. The tail has twice
the precision of the result, and the equation `res * y == x - tail` always holds:

```go
res, tail := dec.Milli.FromUInt64(2).QuoTail(dec.Milli.FromUInt64(3))
// res  == 0.666  (precision 3)
// tail == 0.002  (precision 6)
// 0.666 * 3 == 2 - 0.002 == 1.998
```

They differ in negative-number behavior: `QuoTail` truncates towards zero
(the tail has the sign of the dividend), while `DivTail` uses Euclidean
division (the tail is always non-negative).

### Rescaling

```go
n := dec.Milli.MustParse("123.456")

// Rescale returns a value with a new precision (the extra digits are dropped)
n.Rescale(dec.Deci) // 123.4

// RescaleRem also returns the dropped part with the original precision
rescaled, remainder := n.RescaleRem(dec.Deci) // 123.4 and 0.056
```

### Rounding

```go
n := dec.Nano.MustParse("1.123456789")
n.Round(4, dec.HalfUp)       // 1.1235
n.Round(3, dec.ToZero)       // 1.123
n.Round(3, dec.AwayFromZero) // 1.124

dec.Milli.MustParse("1.5").Round(0, dec.HalfUp)    // 2
dec.Milli.MustParse("-1.5").Round(0, dec.HalfUp)   // -1
dec.Milli.MustParse("-1.5").Round(0, dec.HalfDown) // -2
dec.Milli.MustParse("1.5").Round(0, dec.HalfEven)  // 2
```

### Comparison and sign

```go
a, b := dec.Milli.MustParse("1.5"), dec.Nano.MustParse("1.5")
a.Cmp(b)                        // 0 — precisions are coerced before comparison
dec.Nano.MustParse("-2").Sign() // -1
dec.Nano.MustParse("-2").Abs()  // 2
dec.Nano.MustParse("2").Neg()   // -2
```

### Check the bit size of a value with Fit

`Fit` wraps a `Decimal` and checks that its raw units fit into a given size —
useful, for example, before storing a value into a fixed-size on-chain integer:

```go
value := dec.Nano.MustParse("123.456")

fit, ok := value.Fit(dec.Fit128) // ok == false if the value needs more than 128 bits
fit := value.MustFit(dec.Fit256) // panics if the value does not fit

// operations on Fit values report whether the result still fits
maxValue := dec.Nano.FromUnits(dec.Fit128.MaxValue()).MustFit(dec.Fit128)
sum, ok := maxValue.Add(dec.Nano.Unit().MustFit(dec.Fit128)) // ok == false — overflow
half, ok := maxValue.Div(dec.Nano.FromUInt64(2).MustFit(dec.Fit128)) // ok == true
```

Available sizes: `Fit32`, `Fit64`, `Fit128`, `Fit256`, `Fit512`.

### JSON (and other text) serialization

For every precision constant there is a text-marshaling type (`TextDeci`,
`TextCenti`, ..., `TextQuecto`) that encodes the value as a quoted decimal
string and parses it back with the type's precision:

```go
type Balance struct {
    Available *dec.TextNano `json:"available"`
    Locked    *dec.TextNano `json:"locked"`
}

var b Balance
_ = json.Unmarshal([]byte(`{"available": "0.1", "locked": "0.2"}`), &b)
b.Available.GetDecimal().String() // "0.1"

data, err := json.Marshal(Balance{
    Available: dec.TextNanoRef(dec.Nano.MustParse("0.1")),
    Locked:    dec.TextNanoRef(dec.Nano.MustParse("0.2")),
}) // {"available":"0.1","locked":"0.2"}
```

Marshaling checks that the value's precision matches the type's precision and
fails with `ErrTypePrecisionMismatch` otherwise.

## Benchmarks

The `benchmarks` directory contains a comparison with `big.Rat` for both the
immutable and the mutable APIs.
