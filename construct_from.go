package dec

import (
	"math/big"
	"strings"
)

// FromUnits creates Decimal from a raw *big.Int value and a precision
func FromUnits(val *big.Int, precision Precision) (Decimal, error) {
	if precision < 0 || precision >= 128 {
		return Decimal{}, ErrInvalidPrecision
	}
	if !checkNumberSize(val) {
		return Decimal{}, ErrTooBigNumber
	}
	return Decimal{
		precision: precision,
		value:     (&big.Int{}).Set(val),
	}, nil
}

// MustFromUnits the same as FromUnits but panics on error
func MustFromUnits(val *big.Int, precision Precision) Decimal {
	return must(FromUnits(val, precision))
}

// FromUInt64 creates Decimal from a raw uint64 value and a precision
func FromUInt64(val uint64, precision Precision) (Decimal, error) {
	if precision < 0 || precision >= 128 {
		return Decimal{}, ErrInvalidPrecision
	}
	return FromUnits((&big.Int{}).SetUint64(val), precision)
}

// MustFromUInt64 the same as FromUInt64 but panics on error
func MustFromUInt64(val uint64, precision Precision) Decimal {
	return must(FromUInt64(val, precision))
}

// FromInt64 creates Decimal from a raw int64 value and a precision
func FromInt64(val int64, precision Precision) (Decimal, error) {
	if precision < 0 || precision >= 128 {
		return Decimal{}, ErrInvalidPrecision
	}
	return FromUnits((&big.Int{}).SetInt64(val), precision)
}

// MustFromInt64 the same as FromInt64 but panics on error
func MustFromInt64(val int64, precision Precision) Decimal {
	return must(FromInt64(val, precision))
}

// FromZUInt64 creates Decimal using uint64 as an integer part of the value
func FromZUInt64(val uint64, precision Precision) (Decimal, error) {
	if precision < 0 || precision >= 128 {
		return Decimal{}, ErrInvalidPrecision
	}
	value := (&big.Int{}).SetUint64(val)
	value.Mul(value, precision.Multiplier())
	return Decimal{value: value, precision: precision}, nil
}

// MustFromZUInt64 the same as FromZUInt64 but panics on error
func MustFromZUInt64(val uint64, precision Precision) Decimal {
	return must(FromZUInt64(val, precision))
}

// FromZInt64 creates Decimal using int64 as an integer part of the value
func FromZInt64(val int64, precision Precision) (Decimal, error) {
	if precision < 0 || precision >= 128 {
		return Decimal{}, ErrInvalidPrecision
	}
	value := (&big.Int{}).SetInt64(val)
	value.Mul(value, precision.Multiplier())
	return Decimal{value: value, precision: precision}, nil
}

// MustFromZInt64 the same as FromZInt64 but panics on error
func MustFromZInt64(val int64, precision Precision) Decimal {
	return must(FromZInt64(val, precision))
}

// Parse parses decimal number
func Parse(val string, precision Precision) (Decimal, error) {
	if precision < 0 || precision >= 128 {
		return Decimal{}, ErrInvalidPrecision
	}

	s := strings.SplitN(val, ".", 2)

	if len(s) == 0 {
		return Decimal{}, ErrInvalid
	}

	hi, ok := (&big.Int{}).SetString(s[0], 10)
	if !ok {
		return Decimal{}, ErrInvalid
	}

	hi = hi.Mul(hi, (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil))

	if len(s) == 2 {
		loStr := s[1]
		// lo can have max {decimals} digits
		if len(loStr) > int(precision) {
			loStr = loStr[:precision]
		}

		leadZeroes := 0
		for _, sym := range loStr {
			if sym != '0' {
				break
			}
			leadZeroes++
		}

		lo, ok := (&big.Int{}).SetString(loStr, 10)
		if !ok {
			return Decimal{}, ErrInvalid
		}

		digits := len(lo.String()) // =_=
		lo = lo.Mul(lo, (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64((int(precision)-leadZeroes)-digits)), nil))

		hi = hi.Add(hi, lo)
	}

	return FromUnits(hi, precision)
}

// MustParse the same as ParsePrecise but panics on error
func MustParse(val string, precision Precision) Decimal {
	return must(Parse(val, precision))
}

// ParseUnits parse a string of whole number containing integer and fractional part of the value
func ParseUnits(val string, precision Precision) (Decimal, error) {
	if bn, ok := (&big.Int{}).SetString(val, 10); ok {
		return FromUnits(bn, precision)
	}
	return Decimal{}, ErrInvalid
}

// MustParseUnits the same as ParseUnits but panics on error
func MustParseUnits(val string, precision Precision) Decimal {
	return must(ParseUnits(val, precision))
}

func checkNumberSize(val *big.Int) bool {
	if uint((val.BitLen()+7)>>3) > 32 {
		// (bitlen + 7)>>3 - the same as ceil(bitlen/8)
		// There are no systems need numbers loger then 32 bytes (256 bits)
		// especially for money
		return false
	}
	return true
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
