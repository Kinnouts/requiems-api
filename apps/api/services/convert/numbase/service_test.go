package numbase

import (
	"errors"
	"testing"
)

func TestService_Convert(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name     string
		value    string
		fromBase int
		toBase   int
		want     string
		wantErr  error
	}{
		{
			name:     "decimal to hex",
			value:    "255",
			fromBase: 10,
			toBase:   16,
			want:     "ff",
		},
		{
			name:     "hex to decimal",
			value:    "ff",
			fromBase: 16,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "hex with 0x prefix to decimal",
			value:    "0xff",
			fromBase: 16,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "decimal to binary",
			value:    "255",
			fromBase: 10,
			toBase:   2,
			want:     "11111111",
		},
		{
			name:     "binary to decimal",
			value:    "11111111",
			fromBase: 2,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "binary with 0b prefix to decimal",
			value:    "0b11111111",
			fromBase: 2,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "decimal to octal",
			value:    "255",
			fromBase: 10,
			toBase:   8,
			want:     "377",
		},
		{
			name:     "octal to decimal",
			value:    "377",
			fromBase: 8,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "octal with 0o prefix to decimal",
			value:    "0o377",
			fromBase: 8,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "hex to binary",
			value:    "ff",
			fromBase: 16,
			toBase:   2,
			want:     "11111111",
		},
		{
			name:     "zero",
			value:    "0",
			fromBase: 10,
			toBase:   16,
			want:     "0",
		},
		{
			name:     "negative decimal to hex",
			value:    "-255",
			fromBase: 10,
			toBase:   16,
			want:     "-ff",
		},
		{
			name:     "same base",
			value:    "42",
			fromBase: 10,
			toBase:   10,
			want:     "42",
		},
		{
			name:     "invalid from base",
			value:    "255",
			fromBase: 3,
			toBase:   10,
			wantErr:  ErrInvalidBase,
		},
		{
			name:     "invalid to base",
			value:    "255",
			fromBase: 10,
			toBase:   5,
			wantErr:  ErrInvalidBase,
		},
		{
			name:     "invalid value for base",
			value:    "xyz",
			fromBase: 10,
			toBase:   16,
			wantErr:  ErrInvalidValue,
		},
		{
			name:     "binary value with decimal digits",
			value:    "29",
			fromBase: 2,
			toBase:   10,
			wantErr:  ErrInvalidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Convert(tt.value, tt.fromBase, tt.toBase)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Result != tt.want {
				t.Errorf("Convert(%q, %d, %d) = %q, want %q", tt.value, tt.fromBase, tt.toBase, got.Result, tt.want)
			}
			if got.Input != tt.value {
				t.Errorf("Convert(%q, %d, %d): input = %q, want %q", tt.value, tt.fromBase, tt.toBase, got.Input, tt.value)
			}
			if got.From != tt.fromBase {
				t.Errorf("Convert(%q, %d, %d): from = %d, want %d", tt.value, tt.fromBase, tt.toBase, got.From, tt.fromBase)
			}
			if got.To != tt.toBase {
				t.Errorf("Convert(%q, %d, %d): to = %d, want %d", tt.value, tt.fromBase, tt.toBase, got.To, tt.toBase)
			}
		})
	}
}
