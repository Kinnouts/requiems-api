package main

import "testing"

func TestNormalise_ValidCountryCode(t *testing.T) {
	t.Parallel()

	r := RawInflationRecord{
		CountryCode: " us ",
		CountryName: " United States ",
		Year:        2020,
		Rate:        2.123456789,
	}

	got := normalise(r)

	if got.CountryCode != "US" {
		t.Errorf("CountryCode = %q, want %q", got.CountryCode, "US")
	}
	if got.CountryName != "United States" {
		t.Errorf("CountryName = %q, want %q", got.CountryName, "United States")
	}
	// Rate is rounded to 4 decimal places.
	if got.Rate != 2.1235 {
		t.Errorf("Rate = %v, want 2.1235", got.Rate)
	}
}

func TestNormalise_RegionalAggregateCleared(t *testing.T) {
	t.Parallel()

	// World Bank uses codes like "EAP" (East Asia Pacific), "EMU", "1A" for
	// regional aggregates. The normalise function clears codes whose length is
	// not exactly 2 characters. Note: 2-character codes like "1W" are NOT
	// cleared (only codes with length ≠ 2 are discarded).
	tests := []struct {
		code    string
		cleared bool
	}{
		{"EAP", true}, // 3 chars
		{"EMU", true}, // 3 chars
		{"", true},    // empty
		{"1W", false}, // 2 chars — retained (numeric chars allowed)
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			t.Parallel()
			r := RawInflationRecord{CountryCode: tt.code, CountryName: "Region"}
			got := normalise(r)
			if tt.cleared && got.CountryCode != "" {
				t.Errorf("CountryCode = %q, want empty for non-2-letter code %q", got.CountryCode, tt.code)
			}
			if !tt.cleared && got.CountryCode == "" {
				t.Errorf("CountryCode was cleared unexpectedly for code %q", tt.code)
			}
		})
	}
}
