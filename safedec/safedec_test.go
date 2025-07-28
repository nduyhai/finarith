package safedec

import (
	"errors"
	"testing"

	finerrors "github.com/nduyhai/finarith/errors"
	"github.com/nduyhai/finarith/rounding"
	"github.com/shopspring/decimal"
)

func TestNewFromString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid decimal",
			value:   "10.50",
			want:    "10.5",
			wantErr: false,
		},
		{
			name:    "integer",
			value:   "100",
			want:    "100",
			wantErr: false,
		},
		{
			name:    "negative decimal",
			value:   "-25.75",
			want:    "-25.75",
			wantErr: false,
		},
		{
			name:    "zero",
			value:   "0",
			want:    "0",
			wantErr: false,
		},
		{
			name:    "invalid decimal",
			value:   "abc",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("NewFromString() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestNewFromFloat(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{
			name:  "positive decimal",
			value: 10.5,
			want:  "10.5",
		},
		{
			name:  "negative decimal",
			value: -25.75,
			want:  "-25.75",
		},
		{
			name:  "integer",
			value: 100.0,
			want:  "100",
		},
		{
			name:  "zero",
			value: 0.0,
			want:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromFloat(tt.value)
			if got.String() != tt.want {
				t.Errorf("NewFromFloat() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestNewFromInt(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  string
	}{
		{
			name:  "positive integer",
			value: 100,
			want:  "100",
		},
		{
			name:  "negative integer",
			value: -25,
			want:  "-25",
		},
		{
			name:  "zero",
			value: 0,
			want:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromInt(tt.value)
			if got.String() != tt.want {
				t.Errorf("NewFromInt() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestDecimal_Value(t *testing.T) {
	original := decimal.NewFromFloat(10.5)
	d := New(original)

	if !d.Value().Equal(original) {
		t.Errorf("Value() = %v, want %v", d.Value(), original)
	}
}

func TestDecimal_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "decimal",
			value: "10.5",
			want:  "10.5",
		},
		{
			name:  "integer",
			value: "100",
			want:  "100",
		},
		{
			name:  "negative",
			value: "-25.75",
			want:  "-25.75",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := NewFromString(tt.value)
			if d.String() != tt.want {
				t.Errorf("String() = %v, want %v", d.String(), tt.want)
			}
		})
	}
}

func TestDecimal_Mul(t *testing.T) {
	tests := []struct {
		name   string
		value1 string
		value2 string
		want   string
	}{
		{
			name:   "simple multiplication",
			value1: "10.50",
			value2: "2",
			want:   "21",
		},
		{
			name:   "decimal multiplication",
			value1: "19.99",
			value2: "3",
			want:   "59.97",
		},
		{
			name:   "negative multiplication",
			value1: "-10.5",
			value2: "2",
			want:   "-21",
		},
		{
			name:   "zero multiplication",
			value1: "10.5",
			value2: "0",
			want:   "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1, _ := NewFromString(tt.value1)
			d2, _ := NewFromString(tt.value2)
			result := d1.Mul(d2)
			if result.String() != tt.want {
				t.Errorf("Mul() = %v, want %v", result.String(), tt.want)
			}
		})
	}
}

func TestDecimal_Round(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		places  int32
		mode    rounding.Mode
		want    string
		wantErr bool
	}{
		{
			name:    "round half up",
			value:   "10.555",
			places:  2,
			mode:    rounding.RoundHalfUp,
			want:    "10.56",
			wantErr: false,
		},
		{
			name:    "round half even",
			value:   "10.555",
			places:  2,
			mode:    rounding.RoundHalfEven,
			want:    "10.56",
			wantErr: false,
		},
		{
			name:    "round down",
			value:   "10.555",
			places:  2,
			mode:    rounding.RoundDown,
			want:    "10.55",
			wantErr: false,
		},
		{
			name:    "round up",
			value:   "10.551",
			places:  2,
			mode:    rounding.RoundUp,
			want:    "10.56",
			wantErr: false,
		},
		{
			name:    "invalid rounding mode",
			value:   "10.555",
			places:  2,
			mode:    rounding.Mode(99),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := NewFromString(tt.value)
			result, err := d.Round(tt.places, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Round() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.want {
				t.Errorf("Round() = %v, want %v", result.String(), tt.want)
			}
		})
	}
}

func TestDecimal_DivRound(t *testing.T) {
	tests := []struct {
		name    string
		value1  string
		value2  string
		places  int32
		mode    rounding.Mode
		want    string
		wantErr bool
	}{
		{
			name:    "simple division",
			value1:  "10",
			value2:  "2",
			places:  2,
			mode:    rounding.RoundHalfUp,
			want:    "5",
			wantErr: false,
		},
		{
			name:    "division with rounding",
			value1:  "10",
			value2:  "3",
			places:  2,
			mode:    rounding.RoundHalfUp,
			want:    "3.33",
			wantErr: false,
		},
		{
			name:    "division by zero",
			value1:  "10",
			value2:  "0",
			places:  2,
			mode:    rounding.RoundHalfUp,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid rounding mode",
			value1:  "10",
			value2:  "3",
			places:  2,
			mode:    rounding.Mode(99),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1, _ := NewFromString(tt.value1)
			d2, _ := NewFromString(tt.value2)
			result, err := d1.DivRound(d2, tt.places, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DivRound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.want {
				t.Errorf("DivRound() = %v, want %v", result.String(), tt.want)
			}
			if err != nil && tt.value2 == "0" && !errors.Is(err, finerrors.ErrDivideByZero) {
				t.Errorf("DivRound() error is not ErrDivideByZero: %v", err)
			}
		})
	}
}

func TestDecimal_SubNonNegative(t *testing.T) {
	tests := []struct {
		name    string
		value1  string
		value2  string
		want    string
		wantErr bool
	}{
		{
			name:    "positive result",
			value1:  "50",
			value2:  "30",
			want:    "20",
			wantErr: false,
		},
		{
			name:    "zero result",
			value1:  "50",
			value2:  "50",
			want:    "0",
			wantErr: false,
		},
		{
			name:    "negative result",
			value1:  "50",
			value2:  "60",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1, _ := NewFromString(tt.value1)
			d2, _ := NewFromString(tt.value2)
			result, err := d1.SubNonNegative(d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubNonNegative() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.want {
				t.Errorf("SubNonNegative() = %v, want %v", result.String(), tt.want)
			}
			if err != nil && !errors.Is(err, finerrors.ErrNegativeValue) {
				t.Errorf("SubNonNegative() error is not ErrNegativeValue: %v", err)
			}
		})
	}
}

func TestZero(t *testing.T) {
	zero := Zero()
	if !zero.IsZero() {
		t.Errorf("Zero() = %v, want 0", zero.String())
	}
}

func TestOne(t *testing.T) {
	one := One()
	if one.String() != "1" {
		t.Errorf("One() = %v, want 1", one.String())
	}
}

func TestMinValue(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{
			name: "a < b",
			a:    "10",
			b:    "20",
			want: "10",
		},
		{
			name: "a > b",
			a:    "30",
			b:    "20",
			want: "20",
		},
		{
			name: "a = b",
			a:    "20",
			b:    "20",
			want: "20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewFromString(tt.a)
			b, _ := NewFromString(tt.b)
			result := MinValue(a, b)
			if result.String() != tt.want {
				t.Errorf("MinValue() = %v, want %v", result.String(), tt.want)
			}
		})
	}
}

func TestMaxValue(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{
			name: "a < b",
			a:    "10",
			b:    "20",
			want: "20",
		},
		{
			name: "a > b",
			a:    "30",
			b:    "20",
			want: "30",
		},
		{
			name: "a = b",
			a:    "20",
			b:    "20",
			want: "20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewFromString(tt.a)
			b, _ := NewFromString(tt.b)
			result := MaxValue(a, b)
			if result.String() != tt.want {
				t.Errorf("MaxValue() = %v, want %v", result.String(), tt.want)
			}
		})
	}
}
