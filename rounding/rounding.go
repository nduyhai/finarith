// Package rounding provides configurable rounding strategies for financial calculations.
package rounding

import (
	"math"

	"github.com/nduyhai/finarith/errors"
)

// Mode represents a rounding mode.
type Mode int

// Rounding modes
const (
	// RoundDown rounds toward zero (truncate).
	RoundDown Mode = iota

	// RoundUp rounds away from zero.
	RoundUp

	// RoundHalfUp rounds to nearest, with ties away from zero.
	// This is the most common rounding mode in financial calculations.
	RoundHalfUp

	// RoundHalfDown rounds to nearest, with ties toward zero.
	RoundHalfDown

	// RoundHalfEven rounds to nearest, with ties to even (banker's rounding).
	// This is the default IEEE 754 rounding mode.
	RoundHalfEven

	// RoundCeiling rounds toward positive infinity.
	RoundCeiling

	// RoundFloor rounds toward negative infinity.
	RoundFloor
)

// String returns the string representation of the rounding mode.
func (m Mode) String() string {
	switch m {
	case RoundDown:
		return "round_down"
	case RoundUp:
		return "round_up"
	case RoundHalfUp:
		return "round_half_up"
	case RoundHalfDown:
		return "round_half_down"
	case RoundHalfEven:
		return "round_half_even"
	case RoundCeiling:
		return "round_ceiling"
	case RoundFloor:
		return "round_floor"
	default:
		return "unknown"
	}
}

// RoundFloat64 rounds a float64 value to the specified number of decimal places
// using the specified rounding mode.
func RoundFloat64(value float64, decimals int, mode Mode) (float64, error) {
	if decimals < 0 {
		return 0, errors.ErrInvalidPrecision
	}

	// Special cases
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return value, nil
	}

	// Calculate the multiplier for the specified number of decimal places
	multiplier := math.Pow10(decimals)

	// Apply the appropriate rounding mode
	switch mode {
	case RoundDown:
		return math.Trunc(value*multiplier) / multiplier, nil

	case RoundUp:
		if value >= 0 {
			return math.Ceil(value*multiplier) / multiplier, nil
		}
		return math.Floor(value*multiplier) / multiplier, nil

	case RoundHalfUp:
		if value >= 0 {
			return math.Floor(value*multiplier+0.5) / multiplier, nil
		}
		return math.Ceil(value*multiplier-0.5) / multiplier, nil

	case RoundHalfDown:
		if value >= 0 {
			return math.Floor(value*multiplier+0.5-1e-10) / multiplier, nil
		}
		return math.Ceil(value*multiplier-0.5+1e-10) / multiplier, nil

	case RoundHalfEven:
		// Multiply by the scaling factor
		scaled := value * multiplier

		// Get the integer and fractional parts
		intPart, fracPart := math.Modf(scaled)

		// Check if we're exactly at the halfway point
		if math.Abs(fracPart) == 0.5 {
			// Round to even
			if math.Mod(intPart, 2) == 0 {
				// Even, round down
				return intPart / multiplier, nil
			} else {
				// Odd, round up
				if intPart >= 0 {
					return (intPart + 1) / multiplier, nil
				}
				return (intPart - 1) / multiplier, nil
			}
		}

		// Not at halfway point, use regular rounding
		return math.Round(scaled) / multiplier, nil

	case RoundCeiling:
		return math.Ceil(value*multiplier) / multiplier, nil

	case RoundFloor:
		return math.Floor(value*multiplier) / multiplier, nil

	default:
		return 0, errors.ErrInvalidRounding
	}
}

// RoundInt64 rounds an int64 value to the nearest multiple of the specified unit
// using the specified rounding mode.
func RoundInt64(value, unit int64, mode Mode) (int64, error) {
	if unit <= 0 {
		return 0, errors.ErrInvalidPrecision
	}

	// Apply the appropriate rounding mode
	switch mode {
	case RoundDown:
		// Truncate toward zero
		if value >= 0 {
			return (value / unit) * unit, nil
		}
		return ((value + unit - 1) / unit) * unit, nil

	case RoundUp:
		// Round away from zero
		if value >= 0 {
			return ((value + unit - 1) / unit) * unit, nil
		}
		return (value / unit) * unit, nil

	case RoundHalfUp:
		// Round to nearest, ties away from zero
		halfUnit := unit / 2
		if value >= 0 {
			return ((value + halfUnit) / unit) * unit, nil
		}
		return ((value - halfUnit) / unit) * unit, nil

	case RoundHalfDown:
		// Round to nearest, ties toward zero
		halfUnit := unit / 2
		if value >= 0 {
			if value%unit == halfUnit {
				return (value / unit) * unit, nil
			}
			return ((value + halfUnit) / unit) * unit, nil
		}
		if -value%unit == halfUnit {
			return (value / unit) * unit, nil
		}
		return ((value - halfUnit) / unit) * unit, nil

	case RoundHalfEven:
		// Round to nearest, ties to even
		halfUnit := unit / 2
		if value%unit == halfUnit || value%unit == -halfUnit {
			// At the halfway point, round to even
			quotient := value / unit
			if quotient%2 == 0 {
				// Even, round down
				return quotient * unit, nil
			}
			// Odd, round up
			if quotient >= 0 {
				return (quotient + 1) * unit, nil
			}
			return (quotient - 1) * unit, nil
		}
		// Not at halfway point, use regular rounding
		if value >= 0 {
			return ((value + halfUnit) / unit) * unit, nil
		}
		return ((value - halfUnit) / unit) * unit, nil

	case RoundCeiling:
		// Round toward positive infinity
		if value%unit == 0 {
			return value, nil
		}
		if value > 0 {
			return ((value + unit - 1) / unit) * unit, nil
		}
		return (value / unit) * unit, nil

	case RoundFloor:
		// Round toward negative infinity
		if value%unit == 0 {
			return value, nil
		}
		if value > 0 {
			return (value / unit) * unit, nil
		}
		return ((value - unit + 1) / unit) * unit, nil

	default:
		return 0, errors.ErrInvalidRounding
	}
}
