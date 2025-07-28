# Finarith

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/finarith)](LICENSE)

A comprehensive financial arithmetic library for Go that provides safe, precise, and domain-aware mathematical operations.

## Features

- ✅ Int64 and Uint64 overflow-safe Add/Sub/Mul operations
- ✅ Finance-friendly wrappers for shopspring/decimal
- ✅ Shared error types: ErrOverflow, ErrDivideByZero, etc.
- ✅ Configurable rounding rules (e.g., round half-up, banker's rounding)
- ✅ Domain-specific rules: "cannot subtract below zero", "max transfer limit"

## Installation

```bash
go get github.com/nduyhai/finarith
```

## Development

### Build Commands

The project uses a Makefile to manage common development tasks:

```bash
# Download dependencies
make deps

# Run tests
make test

# Run linting
make lint

# Clean build artifacts
make clean

# Show all available commands
make help
```

## Usage Examples

### Integer Overflow-Safe Operations

```go
import "github.com/nduyhai/finarith/safeint"

// Safe addition with overflow detection
result, err := safeint.Add(9223372036854775800, 10)
if err != nil {
    // Handle overflow error
}

// Subtraction with floor (cannot go below zero)
result, err = safeint.SubWithFloor(100, 30, 0)
if err != nil {
    // Handle error if result would be below floor
}
```

### Unsigned Integer Overflow-Safe Operations

```go
import "github.com/nduyhai/finarith/safeuint"

// Safe addition with overflow detection
result, err := safeuint.Add(18446744073709551610, 10)
if err != nil {
    // Handle overflow error
}

// Subtraction with underflow detection
result, err = safeuint.Sub(100, 50)
if err != nil {
    // Handle underflow error (would occur if 50 > 100)
}

// Subtraction with floor (cannot go below specified value)
result, err = safeuint.SubWithFloor(100, 30, 50)
if err != nil {
    // Handle error if result would be below floor
}
```

### Decimal Operations with Rounding

```go
import (
    "github.com/nduyhai/finarith/safedec"
    "github.com/nduyhai/finarith/rounding"
)

// Create decimal values
price, _ := safedec.NewFromString("19.99")
quantity, _ := safedec.NewFromString("3")

// Multiply with proper rounding
total := price.Mul(quantity)

// Round to 2 decimal places using different rounding modes
roundedHalfUp, _ := total.Round(2, rounding.RoundHalfUp)
roundedHalfEven, _ := total.Round(2, rounding.RoundHalfEven)

// Division with rounding
unitPrice, _ := safedec.NewFromString("10.00")
units, _ := safedec.NewFromString("3")
result, err := unitPrice.DivRound(units, 2, rounding.RoundHalfUp)
```

### Domain-Specific Rules

```go
import (
    "github.com/nduyhai/finarith/safedec"
    "github.com/nduyhai/finarith/rules"
)

// Create a transfer rule
maxAmount, _ := safedec.NewFromString("1000.00")
minAmount, _ := safedec.NewFromString("10.00")
dailyLimit, _ := safedec.NewFromString("5000.00")
transferRule := rules.NewTransferRule(maxAmount, minAmount, dailyLimit, false)

// Validate a transfer
transferAmount, _ := safedec.NewFromString("500.00")
accountBalance, _ := safedec.NewFromString("600.00")
dailyTotal, _ := safedec.NewFromString("4000.00")

err := transferRule.ValidateTransfer(transferAmount, accountBalance, dailyTotal)
if err != nil {
    // Handle validation error
}
```

## Package Overview

### errors

Provides standardized error types for financial arithmetic operations:

- `ErrOverflow`: For arithmetic overflow errors
- `ErrDivideByZero`: For division by zero errors
- `ErrNegativeValue`: For when negative values are not allowed
- `ErrExceedsLimit`: For when values exceed defined limits
- Custom error types with additional context

### safeint

Provides overflow-safe signed integer (int64) arithmetic operations:

- `Add`: Addition with overflow detection
- `Sub`: Subtraction with overflow detection
- `Mul`: Multiplication with overflow detection
- Domain-specific operations with limits

### safeuint

Provides overflow-safe unsigned integer (uint64) arithmetic operations:

- `Add`: Addition with overflow detection
- `Sub`: Subtraction with underflow detection
- `Mul`: Multiplication with overflow detection
- Domain-specific operations with limits

### safedec

Provides finance-friendly wrappers for the shopspring/decimal package:

- Basic arithmetic operations with proper error handling
- Domain-specific operations like `SubNonNegative` and `AddWithLimit`
- Integration with configurable rounding rules

### rounding

Provides configurable rounding strategies:

- `RoundDown`: Rounds toward zero
- `RoundUp`: Rounds away from zero
- `RoundHalfUp`: Rounds to nearest, ties away from zero
- `RoundHalfDown`: Rounds to nearest, ties toward zero
- `RoundHalfEven`: Banker's rounding (rounds to nearest, ties to even)
- `RoundCeiling`: Rounds toward positive infinity
- `RoundFloor`: Rounds toward negative infinity

### rules

Provides domain-specific rules for financial calculations:

- `TransferRule`: For validating financial transfers
- `PricingRule`: For validating prices
- `DiscountRule`: For calculating discounts
- `TaxRule`: For calculating taxes

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.