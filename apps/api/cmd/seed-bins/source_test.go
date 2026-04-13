package main

import (
	"strings"
	"testing"
)

func TestIsValidBINPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		bin  string
		want bool
	}{
		{"411111", true},      // 6-digit valid
		{"41111111", true},    // 8-digit valid
		{"4111", false},       // too short
		{"4111111111", false}, // 10 digits — neither 6 nor 8
		{"41111A", false},     // non-numeric
		{"", false},           // empty
		{"ABCDEF", false},     // all letters
	}

	for _, tt := range tests {
		t.Run(tt.bin, func(t *testing.T) {
			t.Parallel()
			if got := isValidBINPrefix(tt.bin); got != tt.want {
				t.Fatalf("isValidBINPrefix(%q) = %v, want %v", tt.bin, got, tt.want)
			}
		})
	}
}

func TestParseIannuttall(t *testing.T) {
	t.Parallel()

	// Columns: bin,brand,type,category,issuer,alpha_2,alpha_3,country,latitude,longitude,bank_phone,bank_url
	csv := "bin,brand,type,category,issuer,alpha_2,alpha_3,country,latitude,longitude,bank_phone,bank_url\n" +
		"411111,VISA,Credit,Classic,Some Bank,US,USA,United States,37.09,-95.71,+1-800-123-4567,https://example.com\n" +
		"INVALID,VISA,Credit,Classic,Skip Me,US,USA,United States,0,0,,\n"

	records, err := parseIannuttall(strings.NewReader(csv), "test", 0.75)
	if err != nil {
		t.Fatalf("parseIannuttall: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]
	if r.BINPrefix != "411111" {
		t.Errorf("BINPrefix = %q, want %q", r.BINPrefix, "411111")
	}
	if r.Scheme != "VISA" {
		t.Errorf("Scheme = %q, want %q", r.Scheme, "VISA")
	}
	if r.CountryCode != "US" {
		t.Errorf("CountryCode = %q, want %q", r.CountryCode, "US")
	}
	if r.Source != "test" {
		t.Errorf("Source = %q, want %q", r.Source, "test")
	}
	if r.Confidence != 0.75 {
		t.Errorf("Confidence = %v, want %v", r.Confidence, 0.75)
	}
}

func TestParseVenelinkochev(t *testing.T) {
	t.Parallel()

	// Columns: BIN,Brand,Type,Category,Issuer,IssuerPhone,IssuerUrl,isoCode2,isoCode3,CountryName
	csv := "BIN,Brand,Type,Category,Issuer,IssuerPhone,IssuerUrl,isoCode2,isoCode3,CountryName\n" +
		"411111,Visa,Credit,Classic,My Bank,+1555000,https://mybank.com,US,USA,United States\n" +
		"BADINP,Visa,Credit,Classic,Skip,,,US,USA,United States\n"

	records, err := parseVenelinkochev(strings.NewReader(csv), "venelinkochev", 0.80)
	if err != nil {
		t.Fatalf("parseVenelinkochev: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]
	if r.BINPrefix != "411111" {
		t.Errorf("BINPrefix = %q, want %q", r.BINPrefix, "411111")
	}
	if r.IssuerName != "My Bank" {
		t.Errorf("IssuerName = %q, want %q", r.IssuerName, "My Bank")
	}
	if r.CountryName != "United States" {
		t.Errorf("CountryName = %q, want %q", r.CountryName, "United States")
	}
}
