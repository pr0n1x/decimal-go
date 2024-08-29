# Decimal numbers with fixed point

## Features

- Numbers with different precisions can operate with each other
- Library splits mutable and immutable arithmetic API
- Under the hood two numbers: big.Int and precision (decimal places after the point)
- Library has a type `Fit` which allows to check bit-size of a value / operation result

## Examples
### Create 1 (one) with precision of 9 decimal places

```go
// from big.Int
one := dec.FromUnits(new(big.Int).SetUint64(1_000_000_000), dec.Nano)
one := dec.Nano.FromUnits(new(big.Int).SetUint64(1_000_000_000))

// from uint64
one := dec.FromUInt64Units(1_000_000_000), 9)
one := dec.FromUInt64Units(1_000_000_000), dec.Nano)
one := dec.FromInt64Units(1_000_000_000), dec.Nano)
one := dec.FromUInt64(1), dec.Nano)
one := dec.FromInt64(1), dec.Nano)
one := dec.Nano.FromInt64Units(1_000_000_000))
one := dec.Nano.FromInt64(1))

// quite often zero precision is also useful
ten := dec.FromUInt64(10, dec.Z)
ten := dec.Z.FromUInt64(10)
// for example if you want to initialize a counter, but a precision of the operation is unknown
zero := dec.Z.Zero()
zero := dec.Zero(dec.Z)
zero := dec.Z.FromUInt64(0)
```

### Immutable arithmetic API
```go
number := dec.Z.FromInt64(123)
hundred := dec.Nano.FromUInt64(100)
three := dec.Z.FromInt64(3)
twenty := number.Mod(hundred).Sub(three)
```

### Mutable arithmetic API
```go
hundred := dec.Nano.FromUInt64(100)
three := dec.Deci.FromInt64(3)
number := dec.Milli.FromInt64(123) // number == 123
number.Var().Mod(hundred)          // number == 23
number.Var().Sub(three)            // number == 20
// precision will be coerced to maximum of all operated numbers i.e. Nano (9)
```
