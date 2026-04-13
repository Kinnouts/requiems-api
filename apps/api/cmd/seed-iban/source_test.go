package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseOptInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"N/A", 0},
		{"-", 0},
		{"0", 0},
		{"4", 4},
		{"32767", 32767},   // max int16
		{"-32768", -32768}, // min int16
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			if got := parseOptInt(tt.input); got != tt.want {
				t.Fatalf("parseOptInt(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestRawIBANCountry_BankAndAccountOffsets(t *testing.T) {
	t.Parallel()

	// Germany: BBAN is 18 chars (bank=0..7, branch=8..12, account=13..17)
	de := RawIBANCountry{
		CountryCode:   "DE",
		IBANLength:    22,
		BBANLength:    18,
		BankIDStart:   0,
		BankIDEnd:     7,
		BranchIDStart: 8,
		BranchIDEnd:   12,
	}

	if de.BankOffset() != 0 {
		t.Errorf("BankOffset = %d, want 0", de.BankOffset())
	}
	if de.BankLength() != 8 { // 7-0+1
		t.Errorf("BankLength = %d, want 8", de.BankLength())
	}
	if de.AccountOffset() != 13 { // BranchIDEnd(12)+1
		t.Errorf("AccountOffset = %d, want 13", de.AccountOffset())
	}
	if de.AccountLength() != 5 { // BBANLength(18) - 1 - AccountOffset(13) + 1
		t.Errorf("AccountLength = %d, want 5", de.AccountLength())
	}
}

func TestRawIBANCountry_NoBranchCode(t *testing.T) {
	t.Parallel()

	// A country with no branch code: account starts right after the bank code.
	r := RawIBANCountry{
		BBANLength:    14,
		BankIDStart:   0,
		BankIDEnd:     3,
		BranchIDStart: 0,
		BranchIDEnd:   0, // no distinct branch (start == end, not > BankIDEnd)
	}

	if r.AccountOffset() != 4 { // BankIDEnd(3)+1
		t.Errorf("AccountOffset = %d, want 4", r.AccountOffset())
	}
	if r.AccountLength() != 10 { // BBANLength(14) - 1 - AccountOffset(4) + 1
		t.Errorf("AccountLength = %d, want 10", r.AccountLength())
	}
}

func TestFetchAndParse_LocalMockServer(t *testing.T) {
	t.Parallel()

	// registry.txt format: pipe-separated, first line is header.
	// Columns (0-indexed): 0=code, 1=name, 4=bban_format, 6=bban_len, 10=iban_len,
	//   11=bank_start, 12=bank_end, 13=branch_start, 14=branch_end, 16=sepa
	registry := "code|name|x|x|bban_format|x|bban_len|x|x|x|iban_len|bank_start|bank_end|branch_start|branch_end|x|sepa\n" +
		"DE|Germany|x|x|8!n10!n|x|18|x|x|x|22|0|7|8|12|x|1\n" +
		"US|United States|x|x|9!n|x|9|x|x|x|13|0|3|N/A|N/A|x|0\n"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(registry))
	}))
	defer srv.Close()

	countries, err := fetchAndParse(srv.URL)
	if err != nil {
		t.Fatalf("fetchAndParse: %v", err)
	}
	if len(countries) != 2 {
		t.Fatalf("expected 2 countries, got %d", len(countries))
	}

	de := countries[0]
	if de.CountryCode != "DE" {
		t.Errorf("CountryCode = %q, want DE", de.CountryCode)
	}
	if de.IBANLength != 22 {
		t.Errorf("IBANLength = %d, want 22", de.IBANLength)
	}
	if de.BBANLength != 18 {
		t.Errorf("BBANLength = %d, want 18", de.BBANLength)
	}
	if !de.SEPAMember {
		t.Error("expected DE to be a SEPA member")
	}

	us := countries[1]
	if us.SEPAMember {
		t.Error("expected US not to be a SEPA member")
	}
}

func TestFetchAndParse_NoValidRows(t *testing.T) {
	t.Parallel()

	// Valid header but all data rows are too short — expect an error.
	registry := "code|name|x|x|bban_format|x|bban_len|x|x|x|iban_len|bank_start|bank_end|branch_start|branch_end|x|sepa\n" +
		"DE|Germany\n" // too few fields

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(registry))
	}))
	defer srv.Close()

	_, err := fetchAndParse(srv.URL)
	if err == nil {
		t.Fatal("expected error for registry with no valid rows, got nil")
	}
}

func TestFetchAndParse_HTTPError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := fetchAndParse(srv.URL)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
}
