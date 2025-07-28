package rounding

import (
	"errors"
	"math"
	"testing"

	finerrors "github.com/nduyhai/finarith/errors"
)

func TestMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want string
	}{
		{
			name: "RoundDown",
			mode: RoundDown,
			want: "round_down",
		},
		{
			name: "RoundUp",
			mode: RoundUp,
			want: "round_up",
		},
		{
			name: "RoundHalfUp",
			mode: RoundHalfUp,
			want: "round_half_up",
		},
		{
			name: "RoundHalfDown",
			mode: RoundHalfDown,
			want: "round_half_down",
		},
		{
			name: "RoundHalfEven",
			mode: RoundHalfEven,
			want: "round_half_even",
		},
		{
			name: "RoundCeiling",
			mode: RoundCeiling,
			want: "round_ceiling",
		},
		{
			name: "RoundFloor",
			mode: RoundFloor,
			want: "round_floor",
		},
		{
			name: "Unknown mode",
			mode: Mode(99),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("Mode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundFloat64(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		decimals int
		mode     Mode
		want     float64
		wantErr  bool
	}{
		// RoundDown tests
		{
			name:     "RoundDown positive",
			value:    10.555,
			decimals: 2,
			mode:     RoundDown,
			want:     10.55,
			wantErr:  false,
		},
		{
			name:     "RoundDown negative",
			value:    -10.555,
			decimals: 2,
			mode:     RoundDown,
			want:     -10.55,
			wantErr:  false,
		},
		// RoundUp tests
		{
			name:     "RoundUp positive",
			value:    10.555,
			decimals: 2,
			mode:     RoundUp,
			want:     10.56,
			wantErr:  false,
		},
		{
			name:     "RoundUp negative",
			value:    -10.555,
			decimals: 2,
			mode:     RoundUp,
			want:     -10.56,
			wantErr:  false,
		},
		// RoundHalfUp tests
		{
			name:     "RoundHalfUp positive exactly half",
			value:    10.555,
			decimals: 2,
			mode:     RoundHalfUp,
			want:     10.56,
			wantErr:  false,
		},
		{
			name:     "RoundHalfUp positive less than half",
			value:    10.554,
			decimals: 2,
			mode:     RoundHalfUp,
			want:     10.55,
			wantErr:  false,
		},
		{
			name:     "RoundHalfUp negative exactly half",
			value:    -10.555,
			decimals: 2,
			mode:     RoundHalfUp,
			want:     -10.56,
			wantErr:  false,
		},
		// RoundHalfDown tests
		{
			name:     "RoundHalfDown positive exactly half",
			value:    10.555,
			decimals: 2,
			mode:     RoundHalfDown,
			want:     10.55,
			wantErr:  false,
		},
		{
			name:     "RoundHalfDown positive more than half",
			value:    10.556,
			decimals: 2,
			mode:     RoundHalfDown,
			want:     10.56,
			wantErr:  false,
		},
		// RoundHalfEven tests
		{
			name:     "RoundHalfEven half to even",
			value:    10.545,
			decimals: 2,
			mode:     RoundHalfEven,
			want:     10.54,
			wantErr:  false,
		},
		{
			name:     "RoundHalfEven half to odd",
			value:    10.555,
			decimals: 2,
			mode:     RoundHalfEven,
			want:     10.56,
			wantErr:  false,
		},
		// RoundCeiling tests
		{
			name:     "RoundCeiling positive",
			value:    10.551,
			decimals: 2,
			mode:     RoundCeiling,
			want:     10.56,
			wantErr:  false,
		},
		{
			name:     "RoundCeiling negative",
			value:    -10.551,
			decimals: 2,
			mode:     RoundCeiling,
			want:     -10.55,
			wantErr:  false,
		},
		// RoundFloor tests
		{
			name:     "RoundFloor positive",
			value:    10.559,
			decimals: 2,
			mode:     RoundFloor,
			want:     10.55,
			wantErr:  false,
		},
		{
			name:     "RoundFloor negative",
			value:    -10.551,
			decimals: 2,
			mode:     RoundFloor,
			want:     -10.56,
			wantErr:  false,
		},
		// Special cases
		{
			name:     "NaN",
			value:    math.NaN(),
			decimals: 2,
			mode:     RoundHalfUp,
			want:     math.NaN(),
			wantErr:  false,
		},
		{
			name:     "Infinity",
			value:    math.Inf(1),
			decimals: 2,
			mode:     RoundHalfUp,
			want:     math.Inf(1),
			wantErr:  false,
		},
		{
			name:     "Negative decimals",
			value:    10.555,
			decimals: -1,
			mode:     RoundHalfUp,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "Invalid mode",
			value:    10.555,
			decimals: 2,
			mode:     Mode(99),
			want:     0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RoundFloat64(tt.value, tt.decimals, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if math.IsNaN(tt.want) {
					if !math.IsNaN(got) {
						t.Errorf("RoundFloat64() = %v, want NaN", got)
					}
				} else if math.IsInf(tt.want, 0) {
					if !math.IsInf(got, int(math.Copysign(1, tt.want))) {
						t.Errorf("RoundFloat64() = %v, want %v", got, tt.want)
					}
				} else if got != tt.want {
					t.Errorf("RoundFloat64() = %v, want %v", got, tt.want)
				}
			}
			if err != nil && tt.decimals < 0 && !errors.Is(err, finerrors.ErrInvalidPrecision) {
				t.Errorf("RoundFloat64() error is not ErrInvalidPrecision: %v", err)
			}
			if err != nil && tt.mode == Mode(99) && !errors.Is(err, finerrors.ErrInvalidRounding) {
				t.Errorf("RoundFloat64() error is not ErrInvalidRounding: %v", err)
			}
		})
	}
}

func TestRoundInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		unit    int64
		mode    Mode
		want    int64
		wantErr bool
	}{
		// RoundDown tests
		{
			name:    "RoundDown positive",
			value:   155,
			unit:    10,
			mode:    RoundDown,
			want:    150,
			wantErr: false,
		},
		{
			name:    "RoundDown negative",
			value:   -155,
			unit:    10,
			mode:    RoundDown,
			want:    -140,
			wantErr: false,
		},
		// RoundUp tests
		{
			name:    "RoundUp positive",
			value:   155,
			unit:    10,
			mode:    RoundUp,
			want:    160,
			wantErr: false,
		},
		{
			name:    "RoundUp negative",
			value:   -155,
			unit:    10,
			mode:    RoundUp,
			want:    -150,
			wantErr: false,
		},
		// RoundHalfUp tests
		{
			name:    "RoundHalfUp positive exactly half",
			value:   155,
			unit:    10,
			mode:    RoundHalfUp,
			want:    160,
			wantErr: false,
		},
		{
			name:    "RoundHalfUp positive less than half",
			value:   154,
			unit:    10,
			mode:    RoundHalfUp,
			want:    150,
			wantErr: false,
		},
		{
			name:    "RoundHalfUp negative exactly half",
			value:   -155,
			unit:    10,
			mode:    RoundHalfUp,
			want:    -160,
			wantErr: false,
		},
		// RoundHalfDown tests
		{
			name:    "RoundHalfDown positive exactly half",
			value:   155,
			unit:    10,
			mode:    RoundHalfDown,
			want:    150,
			wantErr: false,
		},
		{
			name:    "RoundHalfDown positive more than half",
			value:   156,
			unit:    10,
			mode:    RoundHalfDown,
			want:    160,
			wantErr: false,
		},
		// RoundHalfEven tests
		{
			name:    "RoundHalfEven half to even",
			value:   150,
			unit:    20,
			mode:    RoundHalfEven,
			want:    160,
			wantErr: false,
		},
		{
			name:    "RoundHalfEven half to odd",
			value:   170,
			unit:    20,
			mode:    RoundHalfEven,
			want:    160,
			wantErr: false,
		},
		// RoundCeiling tests
		{
			name:    "RoundCeiling positive",
			value:   151,
			unit:    10,
			mode:    RoundCeiling,
			want:    160,
			wantErr: false,
		},
		{
			name:    "RoundCeiling negative",
			value:   -151,
			unit:    10,
			mode:    RoundCeiling,
			want:    -150,
			wantErr: false,
		},
		// RoundFloor tests
		{
			name:    "RoundFloor positive",
			value:   159,
			unit:    10,
			mode:    RoundFloor,
			want:    150,
			wantErr: false,
		},
		{
			name:    "RoundFloor negative",
			value:   -151,
			unit:    10,
			mode:    RoundFloor,
			want:    -160,
			wantErr: false,
		},
		// Special cases
		{
			name:    "Zero unit",
			value:   155,
			unit:    0,
			mode:    RoundHalfUp,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Negative unit",
			value:   155,
			unit:    -10,
			mode:    RoundHalfUp,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid mode",
			value:   155,
			unit:    10,
			mode:    Mode(99),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RoundInt64(tt.value, tt.unit, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RoundInt64() = %v, want %v", got, tt.want)
			}
			if err != nil && (tt.unit <= 0) && !errors.Is(err, finerrors.ErrInvalidPrecision) {
				t.Errorf("RoundInt64() error is not ErrInvalidPrecision: %v", err)
			}
			if err != nil && tt.mode == Mode(99) && !errors.Is(err, finerrors.ErrInvalidRounding) {
				t.Errorf("RoundInt64() error is not ErrInvalidRounding: %v", err)
			}
		})
	}
}
