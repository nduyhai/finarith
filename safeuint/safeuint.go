// Package safeuint provides overflow-safe unsigned integer arithmetic operations.
package safeuint

import (
	"math"

	"github.com/nduyhai/finarith/errors"
)

// Add performs the addition of two uint64 values with overflow checking.
// Returns an error if the operation results in an overflow.
func Add(a, b uint64) (uint64, error) {
	// Check for overflow: a + b > MaxUint64
	if b > 0 && a > math.MaxUint64-b {
		return 0, errors.NewOverflowError("+", a, b)
	}

	return a + b, nil
}

// Sub performs subtraction of two uint64 values with underflow checking.
// Returns an error if the operation would result in a negative value (underflow).
func Sub(a, b uint64) (uint64, error) {
	// Check for underflow: a < b
	if a < b {
		return 0, errors.NewOverflowError("-", a, b)
	}

	return a - b, nil
}

// Mul performs multiplication of two uint64 values with overflow checking.
// Returns an error if the operation results in an overflow.
func Mul(a, b uint64) (uint64, error) {
	// Special cases to avoid division by zero in the overflow checks
	if a == 0 || b == 0 {
		return 0, nil
	}

	// Check for overflow: a * b > MaxUint64
	if a > math.MaxUint64/b {
		return 0, errors.NewOverflowError("*", a, b)
	}

	return a * b, nil
}

// AddWithLimit performs addition with a maximum limit check.
// Returns an error if the result exceeds the specified limit.
func AddWithLimit(a, b, limit uint64) (uint64, error) {
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
func SubWithFloor(a, b, floor uint64) (uint64, error) {
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
func MulWithLimit(a, b, limit uint64) (uint64, error) {
	result, err := Mul(a, b)
	if err != nil {
		return 0, err
	}

	if result > limit {
		return 0, errors.NewLimitError(result, limit, "multiplication")
	}

	return result, nil
}