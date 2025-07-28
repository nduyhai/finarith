// Package safedec provides finance-friendly wrappers for the shopspring/decimal package.
package safedec

import (
	"github.com/shopspring/decimal"

	"github.com/nduyhai/finarith/errors"
	"github.com/nduyhai/finarith/rounding"
)

// Decimal represents a fixed-point decimal number with finance-friendly operations.
// It wraps shopspring/decimal.Decimal to provide additional functionality.
type Decimal struct {
	value decimal.Decimal
}

// New creates a new Decimal from a decimal.Decimal value.
func New(value decimal.Decimal) Decimal {
	return Decimal{value: value}
}

// NewFromString creates a new Decimal from a string representation.
func NewFromString(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}
	return Decimal{value: d}, nil
}

// NewFromFloat creates a new Decimal from a float64 value.
func NewFromFloat(value float64) Decimal {
	return Decimal{value: decimal.NewFromFloat(value)}
}

// NewFromInt creates a new Decimal from an int64 value.
func NewFromInt(value int64) Decimal {
	return Decimal{value: decimal.NewFromInt(value)}
}

// Value returns the underlying decimal.Decimal value.
func (d Decimal) Value() decimal.Decimal {
	return d.value
}

// String returns the string representation of the decimal value.
func (d Decimal) String() string {
	return d.value.String()
}

// Float64 returns the float64 representation of the decimal value.
func (d Decimal) Float64() float64 {
	f, _ := d.value.Float64()
	return f
}

// IntPart returns the integer part of the decimal value.
func (d Decimal) IntPart() int64 {
	return d.value.IntPart()
}

// Equal returns true if the decimal values are equal.
func (d Decimal) Equal(other Decimal) bool {
	return d.value.Equal(other.value)
}

// GreaterThan returns true if the decimal value is greater than the other.
func (d Decimal) GreaterThan(other Decimal) bool {
	return d.value.GreaterThan(other.value)
}

// GreaterThanOrEqual returns true if the decimal value is greater than or equal to the other.
func (d Decimal) GreaterThanOrEqual(other Decimal) bool {
	return d.value.GreaterThanOrEqual(other.value)
}

// LessThan returns true if the decimal value is less than the other.
func (d Decimal) LessThan(other Decimal) bool {
	return d.value.LessThan(other.value)
}

// LessThanOrEqual returns true if the decimal value is less than or equal to the other.
func (d Decimal) LessThanOrEqual(other Decimal) bool {
	return d.value.LessThanOrEqual(other.value)
}

// IsZero returns true if the decimal value is zero.
func (d Decimal) IsZero() bool {
	return d.value.IsZero()
}

// IsNegative returns true if the decimal value is negative.
func (d Decimal) IsNegative() bool {
	return d.value.IsNegative()
}

// IsPositive returns true if the decimal value is positive.
func (d Decimal) IsPositive() bool {
	return d.value.IsPositive()
}

// Add adds the decimal values and returns a new Decimal.
func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{value: d.value.Add(other.value)}
}

// AddWithLimit adds the decimal values and returns a new Decimal.
// Returns an error if the result exceeds the specified limit.
func (d Decimal) AddWithLimit(other, limit Decimal) (Decimal, error) {
	result := d.Add(other)
	if result.GreaterThan(limit) {
		return Decimal{}, errors.NewLimitError(result.String(), limit.String(), "addition")
	}
	return result, nil
}

// Sub subtracts the other decimal value from this one and returns a new Decimal.
func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{value: d.value.Sub(other.value)}
}

// SubWithFloor subtracts the other decimal value from this one and returns a new Decimal.
// Returns an error if the result is less than the specified floor.
func (d Decimal) SubWithFloor(other, floor Decimal) (Decimal, error) {
	result := d.Sub(other)
	if result.LessThan(floor) {
		return Decimal{}, errors.NewLimitError(result.String(), floor.String(), "subtraction floor")
	}
	return result, nil
}

// SubNonNegative subtracts the other decimal value from this one and returns a new Decimal.
// Returns an error if the result would be negative.
func (d Decimal) SubNonNegative(other Decimal) (Decimal, error) {
	result := d.Sub(other)
	if result.IsNegative() {
		return Decimal{}, errors.ErrNegativeValue
	}
	return result, nil
}

// Mul multiplies the decimal values and returns a new Decimal.
func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{value: d.value.Mul(other.value)}
}

// MulWithLimit multiplies the decimal values and returns a new Decimal.
// Returns an error if the result exceeds the specified limit.
func (d Decimal) MulWithLimit(other, limit Decimal) (Decimal, error) {
	result := d.Mul(other)
	if result.GreaterThan(limit) {
		return Decimal{}, errors.NewLimitError(result.String(), limit.String(), "multiplication")
	}
	return result, nil
}

// Div divides this decimal value by the other and returns a new Decimal.
// Returns an error if the divisor is zero.
func (d Decimal) Div(other Decimal) (Decimal, error) {
	if other.IsZero() {
		return Decimal{}, errors.ErrDivideByZero
	}
	return Decimal{value: d.value.Div(other.value)}, nil
}

// DivRound divides this decimal value by the other, rounds to the specified number of decimal places
// using the specified rounding mode, and returns a new Decimal.
// Returns an error if the divisor is zero or if the rounding mode is invalid.
func (d Decimal) DivRound(other Decimal, places int32, mode rounding.Mode) (Decimal, error) {
	if other.IsZero() {
		return Decimal{}, errors.ErrDivideByZero
	}
	
	// Perform the division
	result := d.value.Div(other.value)
	
	// Apply the rounding mode
	switch mode {
	case rounding.RoundDown:
		return Decimal{value: result.RoundDown(places)}, nil
	case rounding.RoundUp:
		return Decimal{value: result.RoundUp(places)}, nil
	case rounding.RoundHalfUp:
		return Decimal{value: result.Round(places)}, nil
	case rounding.RoundHalfEven:
		return Decimal{value: result.RoundBank(places)}, nil
	case rounding.RoundCeiling:
		return Decimal{value: result.RoundCeil(places)}, nil
	case rounding.RoundFloor:
		return Decimal{value: result.RoundFloor(places)}, nil
	default:
		return Decimal{}, errors.ErrInvalidRounding
	}
}

// Round rounds the decimal value to the specified number of decimal places
// using the specified rounding mode and returns a new Decimal.
func (d Decimal) Round(places int32, mode rounding.Mode) (Decimal, error) {
	switch mode {
	case rounding.RoundDown:
		return Decimal{value: d.value.RoundDown(places)}, nil
	case rounding.RoundUp:
		return Decimal{value: d.value.RoundUp(places)}, nil
	case rounding.RoundHalfUp:
		return Decimal{value: d.value.Round(places)}, nil
	case rounding.RoundHalfEven:
		return Decimal{value: d.value.RoundBank(places)}, nil
	case rounding.RoundCeiling:
		return Decimal{value: d.value.RoundCeil(places)}, nil
	case rounding.RoundFloor:
		return Decimal{value: d.value.RoundFloor(places)}, nil
	default:
		return Decimal{}, errors.ErrInvalidRounding
	}
}

// Abs returns the absolute value of the decimal as a new Decimal.
func (d Decimal) Abs() Decimal {
	return Decimal{value: d.value.Abs()}
}

// Neg returns the negation of the decimal as a new Decimal.
func (d Decimal) Neg() Decimal {
	return Decimal{value: d.value.Neg()}
}

// Truncate truncates the decimal to the specified number of decimal places.
func (d Decimal) Truncate(places int32) Decimal {
	return Decimal{value: d.value.Truncate(places)}
}

// Zero returns a decimal with value 0.
func Zero() Decimal {
	return Decimal{value: decimal.Zero}
}

// One returns a decimal with value 1.
func One() Decimal {
	return Decimal{value: decimal.NewFromInt(1)}
}

// MinValue returns the minimum of the two decimal values.
func MinValue(a, b Decimal) Decimal {
	if a.LessThan(b) {
		return a
	}
	return b
}

// MaxValue returns the maximum of the two decimal values.
func MaxValue(a, b Decimal) Decimal {
	if a.GreaterThan(b) {
		return a
	}
	return b
}