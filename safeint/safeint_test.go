package safeint

import (
	"errors"
	"math"
	"testing"

	finerrors "github.com/nduyhai/finarith/errors"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr bool
	}{
		{
			name:    "simple addition",
			a:       100,
			b:       200,
			want:    300,
			wantErr: false,
		},
		{
			name:    "zero addition",
			a:       100,
			b:       0,
			want:    100,
			wantErr: false,
		},
		{
			name:    "negative numbers",
			a:       -100,
			b:       -200,
			want:    -300,
			wantErr: false,
		},
		{
			name:    "mixed signs",
			a:       100,
			b:       -50,
			want:    50,
			wantErr: false,
		},
		{
			name:    "large numbers",
			a:       math.MaxInt64 - 10,
			b:       5,
			want:    math.MaxInt64 - 5,
			wantErr: false,
		},
		{
			name:    "positive overflow",
			a:       math.MaxInt64 - 5,
			b:       10,
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative overflow",
			a:       math.MinInt64 + 5,
			b:       -10,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
			if err != nil && !errors.Is(err, finerrors.ErrOverflow) {
				t.Errorf("Add() error is not ErrOverflow: %v", err)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr bool
	}{
		{
			name:    "simple subtraction",
			a:       200,
			b:       100,
			want:    100,
			wantErr: false,
		},
		{
			name:    "zero subtraction",
			a:       100,
			b:       0,
			want:    100,
			wantErr: false,
		},
		{
			name:    "negative result",
			a:       100,
			b:       200,
			want:    -100,
			wantErr: false,
		},
		{
			name:    "negative numbers",
			a:       -100,
			b:       -200,
			want:    100,
			wantErr: false,
		},
		{
			name:    "mixed signs",
			a:       100,
			b:       -50,
			want:    150,
			wantErr: false,
		},
		{
			name:    "positive overflow",
			a:       math.MaxInt64,
			b:       math.MinInt64,
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative overflow",
			a:       math.MinInt64,
			b:       math.MaxInt64,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sub(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
			if err != nil && !errors.Is(err, finerrors.ErrOverflow) {
				t.Errorf("Sub() error is not ErrOverflow: %v", err)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr bool
	}{
		{
			name:    "simple multiplication",
			a:       10,
			b:       20,
			want:    200,
			wantErr: false,
		},
		{
			name:    "zero multiplication",
			a:       10,
			b:       0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "negative numbers",
			a:       -10,
			b:       -20,
			want:    200,
			wantErr: false,
		},
		{
			name:    "mixed signs",
			a:       10,
			b:       -20,
			want:    -200,
			wantErr: false,
		},
		{
			name:    "large numbers",
			a:       math.MaxInt64 / 10,
			b:       5,
			want:    (math.MaxInt64 / 10) * 5,
			wantErr: false,
		},
		{
			name:    "positive overflow",
			a:       math.MaxInt64 / 10,
			b:       11,
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative overflow",
			a:       math.MinInt64 / 10,
			b:       11,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
			if err != nil && !errors.Is(err, finerrors.ErrOverflow) {
				t.Errorf("Mul() error is not ErrOverflow: %v", err)
			}
		})
	}
}

func TestAddWithLimit(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		limit   int64
		want    int64
		wantErr bool
	}{
		{
			name:    "within limit",
			a:       10,
			b:       20,
			limit:   50,
			want:    30,
			wantErr: false,
		},
		{
			name:    "at limit",
			a:       10,
			b:       20,
			limit:   30,
			want:    30,
			wantErr: false,
		},
		{
			name:    "exceeds limit",
			a:       10,
			b:       30,
			limit:   30,
			want:    0,
			wantErr: true,
		},
		{
			name:    "overflow",
			a:       math.MaxInt64,
			b:       1,
			limit:   100,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddWithLimit(tt.a, tt.b, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddWithLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddWithLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubWithFloor(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		floor   int64
		want    int64
		wantErr bool
	}{
		{
			name:    "above floor",
			a:       30,
			b:       10,
			floor:   10,
			want:    20,
			wantErr: false,
		},
		{
			name:    "at floor",
			a:       30,
			b:       20,
			floor:   10,
			want:    10,
			wantErr: false,
		},
		{
			name:    "below floor",
			a:       10,
			b:       20,
			floor:   0,
			want:    0,
			wantErr: true,
		},
		{
			name:    "overflow",
			a:       math.MinInt64,
			b:       1,
			floor:   0,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubWithFloor(tt.a, tt.b, tt.floor)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubWithFloor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SubWithFloor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMulWithLimit(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		limit   int64
		want    int64
		wantErr bool
	}{
		{
			name:    "within limit",
			a:       10,
			b:       5,
			limit:   100,
			want:    50,
			wantErr: false,
		},
		{
			name:    "at limit",
			a:       10,
			b:       5,
			limit:   50,
			want:    50,
			wantErr: false,
		},
		{
			name:    "exceeds limit",
			a:       10,
			b:       10,
			limit:   50,
			want:    0,
			wantErr: true,
		},
		{
			name:    "overflow",
			a:       math.MaxInt64,
			b:       2,
			limit:   100,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulWithLimit(tt.a, tt.b, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("MulWithLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MulWithLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}