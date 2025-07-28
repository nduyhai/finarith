// Package rules provides domain-specific rules for financial calculations.
package rules

import (
	"github.com/nduyhai/finarith/errors"
	"github.com/nduyhai/finarith/rounding"
)

// TransferRule represents a rule for financial transfers.
type TransferRule struct {
	// MaxAmount is the maximum amount allowed for a single transfer.
	MaxAmount safedec.Decimal

	// MinAmount is the minimum amount allowed for a single transfer.
	MinAmount safedec.Decimal

	// DailyLimit is the maximum total amount allowed per day.
	DailyLimit safedec.Decimal

	// AllowNegativeBalance determines if the source account can have a negative balance after the transfer.
	AllowNegativeBalance bool
}

// NewTransferRule creates a new TransferRule with the specified limits.
func NewTransferRule(maxAmount, minAmount, dailyLimit safedec.Decimal, allowNegativeBalance bool) *TransferRule {
	return &TransferRule{
		MaxAmount:            maxAmount,
		MinAmount:            minAmount,
		DailyLimit:           dailyLimit,
		AllowNegativeBalance: allowNegativeBalance,
	}
}

// ValidateTransfer validates a transfer against the rule.
// Returns an error if the transfer violates any of the rules.
func (r *TransferRule) ValidateTransfer(amount, sourceBalance, dailyTotal safedec.Decimal) error {
	// Check if the amount is within the allowed range
	if amount.LessThan(r.MinAmount) {
		return errors.NewLimitError(amount.String(), r.MinAmount.String(), "minimum transfer")
	}

	if amount.GreaterThan(r.MaxAmount) {
		return errors.NewLimitError(amount.String(), r.MaxAmount.String(), "maximum transfer")
	}

	// Check if the daily limit would be exceeded
	newDailyTotal := dailyTotal.Add(amount)
	if newDailyTotal.GreaterThan(r.DailyLimit) {
		return errors.NewLimitError(newDailyTotal.String(), r.DailyLimit.String(), "daily transfer")
	}

	// Check if the source account has sufficient balance
	if !r.AllowNegativeBalance {
		newBalance, err := sourceBalance.SubNonNegative(amount)
		if err != nil {
			return errors.NewLimitError(amount.String(), sourceBalance.String(), "available balance")
		}
		_ = newBalance // Avoid unused variable warning
	}

	return nil
}

// PricingRule represents a rule for pricing calculations.
type PricingRule struct {
	// MinPrice is the minimum price allowed.
	MinPrice safedec.Decimal

	// MaxPrice is the maximum price allowed.
	MaxPrice safedec.Decimal

	// AllowZeroPrice determines if a zero price is allowed.
	AllowZeroPrice bool

	// AllowNegativePrice determines if a negative price is allowed.
	AllowNegativePrice bool
}

// NewPricingRule creates a new PricingRule with the specified constraints.
func NewPricingRule(minPrice, maxPrice safedec.Decimal, allowZeroPrice, allowNegativePrice bool) *PricingRule {
	return &PricingRule{
		MinPrice:           minPrice,
		MaxPrice:           maxPrice,
		AllowZeroPrice:     allowZeroPrice,
		AllowNegativePrice: allowNegativePrice,
	}
}

// ValidatePrice validates a price against the rule.
// Returns an error if the price violates any of the rules.
func (r *PricingRule) ValidatePrice(price safedec.Decimal) error {
	// Check if zero prices are allowed
	if price.IsZero() && !r.AllowZeroPrice {
		return errors.NewLimitError("0", r.MinPrice.String(), "minimum price")
	}

	// Check if a negative price is allowed
	if price.IsNegative() && !r.AllowNegativePrice {
		return errors.ErrNegativeValue
	}

	// Check if the price is within the allowed range
	if !price.IsZero() && !price.IsNegative() && price.LessThan(r.MinPrice) {
		return errors.NewLimitError(price.String(), r.MinPrice.String(), "minimum price")
	}

	if price.GreaterThan(r.MaxPrice) {
		return errors.NewLimitError(price.String(), r.MaxPrice.String(), "maximum price")
	}

	return nil
}

// DiscountRule represents a rule for applying discounts.
type DiscountRule struct {
	// MaxDiscountPercent is the maximum discount percentage allowed.
	MaxDiscountPercent safedec.Decimal

	// MinPurchaseAmount is the minimum purchase amount required for a discount.
	MinPurchaseAmount safedec.Decimal

	// MaxDiscountAmount is the maximum absolute discount amount allowed.
	MaxDiscountAmount safedec.Decimal
}

// NewDiscountRule creates a new DiscountRule with the specified constraints.
func NewDiscountRule(maxDiscountPercent, minPurchaseAmount, maxDiscountAmount safedec.Decimal) *DiscountRule {
	return &DiscountRule{
		MaxDiscountPercent: maxDiscountPercent,
		MinPurchaseAmount:  minPurchaseAmount,
		MaxDiscountAmount:  maxDiscountAmount,
	}
}

// CalculateDiscount calculates the discount amount based on the purchase amount and discount percentage.
// Returns an error if the discount violates any of the rules.
func (r *DiscountRule) CalculateDiscount(purchaseAmount, discountPercent safedec.Decimal) (safedec.Decimal, error) {
	// Check if the purchase amount meets the minimum requirement
	if purchaseAmount.LessThan(r.MinPurchaseAmount) {
		return safedec.Zero(), errors.NewLimitError(purchaseAmount.String(), r.MinPurchaseAmount.String(), "minimum purchase for discount")
	}

	// Check if the discount percentage is within the allowed range
	if discountPercent.IsNegative() {
		return safedec.Zero(), errors.ErrNegativeValue
	}

	if discountPercent.GreaterThan(r.MaxDiscountPercent) {
		return safedec.Zero(), errors.NewLimitError(discountPercent.String(), r.MaxDiscountPercent.String(), "maximum discount percentage")
	}

	// Calculate the discount amount
	hundred := safedec.NewFromInt(100)
	discountAmount, err := purchaseAmount.Mul(discountPercent).Div(hundred)
	if err != nil {
		return safedec.Zero(), err
	}

	// Check if the discount amount exceeds the maximum allowed
	if discountAmount.GreaterThan(r.MaxDiscountAmount) {
		return r.MaxDiscountAmount, nil
	}

	return discountAmount, nil
}

// TaxRule represents a rule for calculating taxes.
type TaxRule struct {
	// TaxRate is the tax rate as a percentage.
	TaxRate safedec.Decimal

	// MinTaxableAmount is the minimum amount that is taxable.
	MinTaxableAmount safedec.Decimal

	// MaxTaxAmount is the maximum tax amount that can be charged.
	MaxTaxAmount safedec.Decimal

	// RoundingMode is the rounding mode to use for tax calculations.
	RoundingMode rounding.Mode

	// RoundingPrecision is the number of decimal places to round to.
	RoundingPrecision int32
}

// NewTaxRule creates a new TaxRule with the specified parameters.
func NewTaxRule(taxRate, minTaxableAmount, maxTaxAmount safedec.Decimal, roundingMode rounding.Mode, roundingPrecision int32) *TaxRule {
	return &TaxRule{
		TaxRate:           taxRate,
		MinTaxableAmount:  minTaxableAmount,
		MaxTaxAmount:      maxTaxAmount,
		RoundingMode:      roundingMode,
		RoundingPrecision: roundingPrecision,
	}
}

// CalculateTax calculates the tax amount based on the taxable amount.
// Returns an error if the tax calculation violates any of the rules.
func (r *TaxRule) CalculateTax(taxableAmount safedec.Decimal) (safedec.Decimal, error) {
	// Check if the amount is taxable
	if taxableAmount.LessThan(r.MinTaxableAmount) {
		return safedec.Zero(), nil
	}

	// Calculate the tax amount
	hundred := safedec.NewFromInt(100)
	taxAmount, err := taxableAmount.Mul(r.TaxRate).Div(hundred)
	if err != nil {
		return safedec.Zero(), err
	}

	// Round the tax amount according to the specified rounding mode and precision
	taxAmount, err = taxAmount.Round(r.RoundingPrecision, r.RoundingMode)
	if err != nil {
		return safedec.Zero(), err
	}

	// Check if the tax amount exceeds the maximum allowed
	if taxAmount.GreaterThan(r.MaxTaxAmount) {
		return r.MaxTaxAmount, nil
	}

	return taxAmount, nil
}
