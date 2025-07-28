// Package safeint provides overflow-safe integer arithmetic operations.
package safeint

import (
	"math"

	"github.com/nduyhai/finarith/errors"
)

// Add performs the addition of two int64 values with overflow checking.
// Returns an error if the operation results in an overflow.
func Add(a, b int64) (int64, error) {
	// Check for positive overflow: a + b > MaxInt64
	if b > 0 && a > math.MaxInt64-b {
		return 0, errors.NewOverflowError("+", a, b)
	}

	// Check for negative overflow: a + b < MinInt64
	if b < 0 && a < math.MinInt64-b {
		return 0, errors.NewOverflowError("+", a, b)
	}

	return a + b, nil
}

// Sub performs subtraction of two int64 values with overflow checking.
// Returns an error if the operation results in an overflow.
func Sub(a, b int64) (int64, error) {
	// Check for positive overflow: a - b > MaxInt64, which can happen when b is very negative
	if b < 0 && a > math.MaxInt64+b {
		return 0, errors.NewOverflowError("-", a, b)
	}

	// Check for negative overflow: a - b < MinInt64, which can happen when b is very positive
	if b > 0 && a < math.MinInt64+b {
		return 0, errors.NewOverflowError("-", a, b)
	}

	return a - b, nil
}

// Mul performs multiplication of two int64 values with overflow checking.
// Returns an error if the operation results in an overflow.
func Mul(a, b int64) (int64, error) {
	// Special cases to avoid division by zero in the overflow checks
	if a == 0 || b == 0 {
		return 0, nil
	}

	// Check for overflow
	if a > 0 && b > 0 {
		// Both positive: check if a > MaxInt64/b
		if a > math.MaxInt64/b {
			return 0, errors.NewOverflowError("*", a, b)
		}
	} else if a < 0 && b < 0 {
		// Both negative: check if a < MaxInt64/b (result will be positive)
		if a < math.MaxInt64/b {
			return 0, errors.NewOverflowError("*", a, b)
		}
	} else if a > 0 && b < 0 {
		// a positive, b negative: check if b < MinInt64/a
		if b < math.MinInt64/a {
			return 0, errors.NewOverflowError("*", a, b)
		}
	} else if a < 0 && b > 0 {
		// a negative, b positive: check if a < MinInt64/b
		if a < math.MinInt64/b {
			return 0, errors.NewOverflowError("*", a, b)
		}
	}

	return a * b, nil
}

// AddWithLimit performs addition with a maximum limit check.
// Returns an error if the result exceeds the specified limit.
func AddWithLimit(a, b, limit int64) (int64, error) {
	result, err := Add(a, b)
	if err != nil {
		return 0, err
	}

	if result > limit {
		return 0, errors.NewLimitError(result, limit, "addition")
	}

	return result, nil
}

// SubWithFloor performs subtraction with a minimum limit (floor) check.
// Returns an error if the result is less than the specified floor.
func SubWithFloor(a, b, floor int64) (int64, error) {
	result, err := Sub(a, b)
	if err != nil {
		return 0, err
	}

	if result < floor {
		return 0, errors.NewLimitError(result, floor, "subtraction floor")
	}

	return result, nil
}

// MulWithLimit performs multiplication with a maximum limit check.
// Returns an error if the result exceeds the specified limit.
func MulWithLimit(a, b, limit int64) (int64, error) {
	result, err := Mul(a, b)
	if err != nil {
		return 0, err
	}

	if result > limit {
		return 0, errors.NewLimitError(result, limit, "multiplication")
	}

	return result, nil
}
