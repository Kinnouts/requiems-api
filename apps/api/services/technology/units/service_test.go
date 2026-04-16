package units

import (
	"errors"
	"testing"
)

func TestService_Convert(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name        string
		from        string
		to          string
		value       float64
		wantResult  float64
		wantFormula string
		wantErr     error
	}{
		// Length
		{
			name: "miles to km",
			from: "miles", to: "km", value: 10,
			wantResult:  16.09344,
			wantFormula: "miles × 1.609344",
		},
		{
			name: "km to miles",
			from: "km", to: "miles", value: 1,
			wantResult:  0.621371,
			wantFormula: "km × 0.621371",
		},
		{
			name: "m to cm",
			from: "m", to: "cm", value: 1,
			wantResult:  100,
			wantFormula: "m × 100",
		},
		{
			name: "ft to in",
			from: "ft", to: "in", value: 1,
			wantResult:  12,
			wantFormula: "ft × 12",
		},
		{
			name: "same unit (m to m)",
			from: "m", to: "m", value: 5,
			wantResult:  5,
			wantFormula: "m × 1",
		},
		// Weight
		{
			name: "kg to lb",
			from: "kg", to: "lb", value: 1,
			wantResult:  2.204624,
			wantFormula: "kg × 2.204624",
		},
		{
			name: "oz to g",
			from: "oz", to: "g", value: 1,
			wantResult:  28.3495,
			wantFormula: "oz × 28.3495",
		},
		// Volume
		{
			name: "l to ml",
			from: "l", to: "ml", value: 1,
			wantResult:  1000,
			wantFormula: "l × 1000",
		},
		{
			name: "gal to l",
			from: "gal", to: "l", value: 1,
			wantResult:  3.78541,
			wantFormula: "gal × 3.78541",
		},
		// Temperature
		{
			name: "celsius to fahrenheit",
			from: "c", to: "f", value: 100,
			wantResult:  212,
			wantFormula: "°C × 9/5 + 32",
		},
		{
			name: "fahrenheit to celsius",
			from: "f", to: "c", value: 32,
			wantResult:  0,
			wantFormula: "(°F − 32) × 5/9",
		},
		{
			name: "celsius to kelvin",
			from: "c", to: "k", value: 0,
			wantResult:  273.15,
			wantFormula: "°C + 273.15",
		},
		{
			name: "kelvin to celsius",
			from: "k", to: "c", value: 273.15,
			wantResult:  0,
			wantFormula: "K − 273.15",
		},
		// Area
		{
			name: "m2 to ft2",
			from: "m2", to: "ft2", value: 1,
			wantResult:  10.763915,
			wantFormula: "m2 × 10.763915",
		},
		// Speed
		{
			name: "mph to km_h",
			from: "mph", to: "km_h", value: 60,
			wantResult:  96.5604,
			wantFormula: "mph × 1.60934",
		},
		// Errors
		{
			name: "unknown from unit",
			from: "furlong", to: "km", value: 1,
			wantErr: ErrUnknownUnit,
		},
		{
			name: "unknown to unit",
			from: "km", to: "parsec", value: 1,
			wantErr: ErrUnknownUnit,
		},
		{
			name: "incompatible units",
			from: "miles", to: "kg", value: 1,
			wantErr: ErrIncompatibleUnits,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Convert(tt.from, tt.to, tt.value)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Result != tt.wantResult {
				t.Errorf("result: got %v, want %v", got.Result, tt.wantResult)
			}

			if got.Formula != tt.wantFormula {
				t.Errorf("formula: got %q, want %q", got.Formula, tt.wantFormula)
			}
		})
	}
}
