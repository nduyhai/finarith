// Package errors provides standardized error types for financial arithmetic operations.
package errors

import (
	"errors"
	"fmt"
)

// Standard errors that can be returned by financial arithmetic operations.
var (
	// ErrOverflow is returned when an arithmetic operation would result in a value that exceeds the
	// representable range.
	ErrOverflow = errors.New("arithmetic overflow")

	// ErrDivideByZero is returned when a division operation has a zero divisor.
	ErrDivideByZero = errors.New("divide by zero")

	// ErrNegativeValue is returned when a negative value is provided where a non-negative value is required.
	ErrNegativeValue = errors.New("negative value not allowed")

	// ErrExceedsLimit is returned when a value exceeds a defined limit.
	ErrExceedsLimit = errors.New("value exceeds limit")

	// ErrInvalidPrecision is returned when an invalid precision is specified.
	ErrInvalidPrecision = errors.New("invalid precision")

	// ErrInvalidRounding is returned when an invalid rounding mode is specified.
	ErrInvalidRounding = errors.New("invalid rounding mode")
)

// OverflowError represents an arithmetic overflow with additional context.
type OverflowError struct {
	Operation string
	A, B      interface{}
}

// Error returns the error message for an OverflowError.
func (e *OverflowError) Error() string {
	return fmt.Sprintf("%s operation would overflow: %v %s %v", e.Operation, e.A, e.Operation, e.B)
}

// Is implements the errors.Is interface.
func (e *OverflowError) Is(target error) bool {
	return target == ErrOverflow
}

// NewOverflowError creates a new OverflowError.
func NewOverflowError(op string, a, b interface{}) *OverflowError {
	return &OverflowError{
		Operation: op,
		A:         a,
		B:         b,
	}
}

// LimitError represents an error when a value exceeds a defined limit.
type LimitError struct {
	Value     interface{}
	Limit     interface{}
	Operation string
}

// Error returns the error message for a LimitError.
func (e *LimitError) Error() string {
	return fmt.Sprintf("%v exceeds %s limit of %v", e.Value, e.Operation, e.Limit)
}

// Is It implements the errors.Is interface.
func (e *LimitError) Is(target error) bool {
	return target == ErrExceedsLimit
}

// NewLimitError creates a new LimitError.
func NewLimitError(value, limit interface{}, operation string) *LimitError {
	return &LimitError{
		Value:     value,
		Limit:     limit,
		Operation: operation,
	}
}
