package assert

import "fmt"

func Ok[T any](v T, ok bool) T {
	if !ok {
		panic("it's not ok")
	}
	return v
}

func Ok2[A any, B any](a A, b B, ok bool) (A, B) {
	if !ok {
		panic("it's not ok")
	}
	return a, b
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must2[A any, B any](a A, b B, err error) (A, B) {
	if err != nil {
		panic(err)
	}
	return a, b
}

func Must3[A any, B any, C any](a A, b B, c C, err error) (A, B, C) {
	if err != nil {
		panic(err)
	}
	return a, b, c
}

func Trust(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			panic(fmt.Errorf("%s: %w", msg[0], err))
		}
	}
}

func NotNil[T any](v *T, msg ...string) *T {
	if v == nil {
		if len(msg) > 0 {
			panic(msg[0])
		}
		panic("value is nil")
	}
	return v
}

func NotEmptySlice[S ~[]T, T any](s S, msg ...string) S {
	if len(s) < 1 {
		if len(msg) > 0 {
			panic(msg[0])
		}
		panic("slice is empty")
	}
	return s
}

func NotEmptyMap[M ~map[K]V, K comparable, V any](m M, msg ...string) M {
	if len(m) < 1 {
		if len(msg) > 0 {
			panic(msg[0])
		}
		panic("map is empty")
	}
	return m
}

func NotZero[N number](n N, msg ...string) N {
	if n == 0 {
		if len(msg) > 0 {
			panic(msg[0])
		}
		panic("value is zero")
	}
	return n
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type float interface {
	~float32 | ~float64
}

type number interface{ integer | float }
