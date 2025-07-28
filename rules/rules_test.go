package rules

import (
	"errors"
	"testing"

	finerrors "github.com/nduyhai/finarith/errors"
	"github.com/nduyhai/finarith/rounding"
	"github.com/nduyhai/finarith/safedec"
)

func TestNewTransferRule(t *testing.T) {
	maxAmount, _ := safedec.NewFromString("1000.00")
	minAmount, _ := safedec.NewFromString("10.00")
	dailyLimit, _ := safedec.NewFromString("5000.00")
	allowNegativeBalance := false

	rule := NewTransferRule(maxAmount, minAmount, dailyLimit, allowNegativeBalance)

	if !rule.MaxAmount.Equal(maxAmount) {
		t.Errorf("NewTransferRule() MaxAmount = %v, want %v", rule.MaxAmount, maxAmount)
	}
	if !rule.MinAmount.Equal(minAmount) {
		t.Errorf("NewTransferRule() MinAmount = %v, want %v", rule.MinAmount, minAmount)
	}
	if !rule.DailyLimit.Equal(dailyLimit) {
		t.Errorf("NewTransferRule() DailyLimit = %v, want %v", rule.DailyLimit, dailyLimit)
	}
	if rule.AllowNegativeBalance != allowNegativeBalance {
		t.Errorf("NewTransferRule() AllowNegativeBalance = %v, want %v", rule.AllowNegativeBalance, allowNegativeBalance)
	}
}

func TestTransferRule_ValidateTransfer(t *testing.T) {
	maxAmount, _ := safedec.NewFromString("1000.00")
	minAmount, _ := safedec.NewFromString("10.00")
	dailyLimit, _ := safedec.NewFromString("5000.00")

	tests := []struct {
		name                 string
		amount               string
		sourceBalance        string
		dailyTotal           string
		allowNegativeBalance bool
		wantErr              bool
		errorType            error
	}{
		{
			name:                 "valid transfer",
			amount:               "500.00",
			sourceBalance:        "600.00",
			dailyTotal:           "4000.00",
			allowNegativeBalance: false,
			wantErr:              false,
		},
		{
			name:                 "below minimum amount",
			amount:               "5.00",
			sourceBalance:        "600.00",
			dailyTotal:           "4000.00",
			allowNegativeBalance: false,
			wantErr:              true,
			errorType:            finerrors.ErrExceedsLimit,
		},
		{
			name:                 "above maximum amount",
			amount:               "1500.00",
			sourceBalance:        "2000.00",
			dailyTotal:           "4000.00",
			allowNegativeBalance: false,
			wantErr:              true,
			errorType:            finerrors.ErrExceedsLimit,
		},
		{
			name:                 "exceeds daily limit",
			amount:               "1000.00",
			sourceBalance:        "2000.00",
			dailyTotal:           "4500.00",
			allowNegativeBalance: false,
			wantErr:              true,
			errorType:            finerrors.ErrExceedsLimit,
		},
		{
			name:                 "insufficient balance",
			amount:               "700.00",
			sourceBalance:        "600.00",
			dailyTotal:           "4000.00",
			allowNegativeBalance: false,
			wantErr:              true,
			errorType:            finerrors.ErrExceedsLimit,
		},
		{
			name:                 "allow negative balance",
			amount:               "700.00",
			sourceBalance:        "600.00",
			dailyTotal:           "4000.00",
			allowNegativeBalance: true,
			wantErr:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewTransferRule(maxAmount, minAmount, dailyLimit, tt.allowNegativeBalance)

			amount, _ := safedec.NewFromString(tt.amount)
			sourceBalance, _ := safedec.NewFromString(tt.sourceBalance)
			dailyTotal, _ := safedec.NewFromString(tt.dailyTotal)

			err := rule.ValidateTransfer(amount, sourceBalance, dailyTotal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorType != nil && !errors.Is(err, tt.errorType) {
				t.Errorf("ValidateTransfer() error type = %v, want %v", err, tt.errorType)
			}
		})
	}
}

func TestNewPricingRule(t *testing.T) {
	minPrice, _ := safedec.NewFromString("10.00")
	maxPrice, _ := safedec.NewFromString("1000.00")
	allowZeroPrice := true
	allowNegativePrice := false

	rule := NewPricingRule(minPrice, maxPrice, allowZeroPrice, allowNegativePrice)

	if !rule.MinPrice.Equal(minPrice) {
		t.Errorf("NewPricingRule() MinPrice = %v, want %v", rule.MinPrice, minPrice)
	}
	if !rule.MaxPrice.Equal(maxPrice) {
		t.Errorf("NewPricingRule() MaxPrice = %v, want %v", rule.MaxPrice, maxPrice)
	}
	if rule.AllowZeroPrice != allowZeroPrice {
		t.Errorf("NewPricingRule() AllowZeroPrice = %v, want %v", rule.AllowZeroPrice, allowZeroPrice)
	}
	if rule.AllowNegativePrice != allowNegativePrice {
		t.Errorf("NewPricingRule() AllowNegativePrice = %v, want %v", rule.AllowNegativePrice, allowNegativePrice)
	}
}

func TestPricingRule_ValidatePrice(t *testing.T) {
	minPrice, _ := safedec.NewFromString("10.00")
	maxPrice, _ := safedec.NewFromString("1000.00")

	tests := []struct {
		name               string
		price              string
		allowZeroPrice     bool
		allowNegativePrice bool
		wantErr            bool
		errorType          error
	}{
		{
			name:               "valid price",
			price:              "500.00",
			allowZeroPrice:     false,
			allowNegativePrice: false,
			wantErr:            false,
		},
		{
			name:               "below minimum price",
			price:              "5.00",
			allowZeroPrice:     false,
			allowNegativePrice: false,
			wantErr:            true,
			errorType:          finerrors.ErrExceedsLimit,
		},
		{
			name:               "above maximum price",
			price:              "1500.00",
			allowZeroPrice:     false,
			allowNegativePrice: false,
			wantErr:            true,
			errorType:          finerrors.ErrExceedsLimit,
		},
		{
			name:               "zero price not allowed",
			price:              "0.00",
			allowZeroPrice:     false,
			allowNegativePrice: false,
			wantErr:            true,
			errorType:          finerrors.ErrExceedsLimit,
		},
		{
			name:               "zero price allowed",
			price:              "0.00",
			allowZeroPrice:     true,
			allowNegativePrice: false,
			wantErr:            false,
		},
		{
			name:               "negative price not allowed",
			price:              "-10.00",
			allowZeroPrice:     true,
			allowNegativePrice: false,
			wantErr:            true,
			errorType:          finerrors.ErrNegativeValue,
		},
		{
			name:               "negative price allowed",
			price:              "-10.00",
			allowZeroPrice:     true,
			allowNegativePrice: true,
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewPricingRule(minPrice, maxPrice, tt.allowZeroPrice, tt.allowNegativePrice)

			price, _ := safedec.NewFromString(tt.price)

			err := rule.ValidatePrice(price)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorType != nil && !errors.Is(err, tt.errorType) {
				t.Errorf("ValidatePrice() error type = %v, want %v", err, tt.errorType)
			}
		})
	}
}

func TestNewDiscountRule(t *testing.T) {
	maxDiscountPercent, _ := safedec.NewFromString("30.00")
	minPurchaseAmount, _ := safedec.NewFromString("100.00")
	maxDiscountAmount, _ := safedec.NewFromString("50.00")

	rule := NewDiscountRule(maxDiscountPercent, minPurchaseAmount, maxDiscountAmount)

	if !rule.MaxDiscountPercent.Equal(maxDiscountPercent) {
		t.Errorf("NewDiscountRule() MaxDiscountPercent = %v, want %v", rule.MaxDiscountPercent, maxDiscountPercent)
	}
	if !rule.MinPurchaseAmount.Equal(minPurchaseAmount) {
		t.Errorf("NewDiscountRule() MinPurchaseAmount = %v, want %v", rule.MinPurchaseAmount, minPurchaseAmount)
	}
	if !rule.MaxDiscountAmount.Equal(maxDiscountAmount) {
		t.Errorf("NewDiscountRule() MaxDiscountAmount = %v, want %v", rule.MaxDiscountAmount, maxDiscountAmount)
	}
}

func TestDiscountRule_CalculateDiscount(t *testing.T) {
	maxDiscountPercent, _ := safedec.NewFromString("30.00")
	minPurchaseAmount, _ := safedec.NewFromString("100.00")
	maxDiscountAmount, _ := safedec.NewFromString("50.00")

	tests := []struct {
		name            string
		purchaseAmount  string
		discountPercent string
		want            string
		wantErr         bool
		errorType       error
	}{
		{
			name:            "valid discount",
			purchaseAmount:  "200.00",
			discountPercent: "15.00",
			want:            "30",
			wantErr:         false,
		},
		{
			name:            "below minimum purchase amount",
			purchaseAmount:  "50.00",
			discountPercent: "15.00",
			want:            "",
			wantErr:         true,
			errorType:       finerrors.ErrExceedsLimit,
		},
		{
			name:            "above maximum discount percent",
			purchaseAmount:  "200.00",
			discountPercent: "40.00",
			want:            "",
			wantErr:         true,
			errorType:       finerrors.ErrExceedsLimit,
		},
		{
			name:            "negative discount percent",
			purchaseAmount:  "200.00",
			discountPercent: "-10.00",
			want:            "",
			wantErr:         true,
			errorType:       finerrors.ErrNegativeValue,
		},
		{
			name:            "exceeds maximum discount amount",
			purchaseAmount:  "500.00",
			discountPercent: "15.00",
			want:            "50",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewDiscountRule(maxDiscountPercent, minPurchaseAmount, maxDiscountAmount)

			purchaseAmount, _ := safedec.NewFromString(tt.purchaseAmount)
			discountPercent, _ := safedec.NewFromString(tt.discountPercent)

			got, err := rule.CalculateDiscount(purchaseAmount, discountPercent)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("CalculateDiscount() = %v, want %v", got.String(), tt.want)
			}

			if err != nil && tt.errorType != nil && !errors.Is(err, tt.errorType) {
				t.Errorf("CalculateDiscount() error type = %v, want %v", err, tt.errorType)
			}
		})
	}
}

func TestNewTaxRule(t *testing.T) {
	taxRate, _ := safedec.NewFromString("10.00")
	minTaxableAmount, _ := safedec.NewFromString("100.00")
	maxTaxAmount, _ := safedec.NewFromString("1000.00")
	roundingMode := rounding.RoundHalfUp
	roundingPrecision := int32(2)

	rule := NewTaxRule(taxRate, minTaxableAmount, maxTaxAmount, roundingMode, roundingPrecision)

	if !rule.TaxRate.Equal(taxRate) {
		t.Errorf("NewTaxRule() TaxRate = %v, want %v", rule.TaxRate, taxRate)
	}
	if !rule.MinTaxableAmount.Equal(minTaxableAmount) {
		t.Errorf("NewTaxRule() MinTaxableAmount = %v, want %v", rule.MinTaxableAmount, minTaxableAmount)
	}
	if !rule.MaxTaxAmount.Equal(maxTaxAmount) {
		t.Errorf("NewTaxRule() MaxTaxAmount = %v, want %v", rule.MaxTaxAmount, maxTaxAmount)
	}
	if rule.RoundingMode != roundingMode {
		t.Errorf("NewTaxRule() RoundingMode = %v, want %v", rule.RoundingMode, roundingMode)
	}
	if rule.RoundingPrecision != roundingPrecision {
		t.Errorf("NewTaxRule() RoundingPrecision = %v, want %v", rule.RoundingPrecision, roundingPrecision)
	}
}

func TestTaxRule_CalculateTax(t *testing.T) {
	taxRate, _ := safedec.NewFromString("10.00")
	minTaxableAmount, _ := safedec.NewFromString("100.00")
	maxTaxAmount, _ := safedec.NewFromString("1000.00")
	roundingMode := rounding.RoundHalfUp
	roundingPrecision := int32(2)

	tests := []struct {
		name          string
		taxableAmount string
		want          string
		wantErr       bool
	}{
		{
			name:          "valid tax calculation",
			taxableAmount: "500.00",
			want:          "50",
			wantErr:       false,
		},
		{
			name:          "below minimum taxable amount",
			taxableAmount: "50.00",
			want:          "0",
			wantErr:       false,
		},
		{
			name:          "exceeds maximum tax amount",
			taxableAmount: "15000.00",
			want:          "1000",
			wantErr:       false,
		},
		{
			name:          "rounding applied",
			taxableAmount: "123.45",
			want:          "12.35",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewTaxRule(taxRate, minTaxableAmount, maxTaxAmount, roundingMode, roundingPrecision)

			taxableAmount, _ := safedec.NewFromString(tt.taxableAmount)

			got, err := rule.CalculateTax(taxableAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("CalculateTax() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
