package swift

import (
	"testing"

	"requiems-api/platform/httpx"
)

func TestSanitizeSWIFT_Valid8Char(t *testing.T) {
	got, err := sanitizeSWIFT("DEUTDEDB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUTDEDBXXX" {
		t.Errorf("expected DEUTDEDBXXX, got %q", got)
	}
}

func TestSanitizeSWIFT_Valid11Char(t *testing.T) {
	got, err := sanitizeSWIFT("DEUTDEDB001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUTDEDB001" {
		t.Errorf("expected DEUTDEDB001, got %q", got)
	}
}

func TestSanitizeSWIFT_Lowercase(t *testing.T) {
	got, err := sanitizeSWIFT("deutdedb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUTDEDBXXX" {
		t.Errorf("expected DEUTDEDBXXX, got %q", got)
	}
}

func TestSanitizeSWIFT_WithSpaces(t *testing.T) {
	got, err := sanitizeSWIFT("  DEUTDEDB  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUTDEDBXXX" {
		t.Errorf("expected DEUTDEDBXXX, got %q", got)
	}
}

func TestSanitizeSWIFT_PrimaryOfficeXXX(t *testing.T) {
	got, err := sanitizeSWIFT("DEUTDEDBXXX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUTDEDBXXX" {
		t.Errorf("expected DEUTDEDBXXX, got %q", got)
	}
}

func TestSanitizeSWIFT_AlphanumericLocation(t *testing.T) {
	// Location code with digit is valid (chars 7-8 are alphanumeric).
	got, err := sanitizeSWIFT("DEUTDE2B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUTDE2BXXX" {
		t.Errorf("expected DEUTDE2BXXX, got %q", got)
	}
}

func TestSanitizeSWIFT_TooShort(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDE")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_TooLong(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB001X")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_9Chars(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB0")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_10Chars(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB01")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_Empty(t *testing.T) {
	_, err := sanitizeSWIFT("")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_DigitInBankCode(t *testing.T) {
	_, err := sanitizeSWIFT("1EUTDEDB")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_DigitInCountryCode(t *testing.T) {
	_, err := sanitizeSWIFT("DEUT1EDB")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_InvalidBranchCode(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB0-1")
	assertAppError(t, err)
}

func TestSanitizeSWIFT_8Char_AppendXXX(t *testing.T) {
	got, err := sanitizeSWIFT("CHASUS33")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 11 {
		t.Errorf("expected 11-char result, got %d: %q", len(got), got)
	}
	if got[8:] != "XXX" {
		t.Errorf("expected branch code XXX, got %q", got[8:])
	}
}

// assertAppError verifies that err is a 400 bad_request *httpx.AppError.
// All sanitizeSWIFT validation errors return this exact status and code.
func assertAppError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	ae, ok := err.(*httpx.AppError)
	if !ok {
		t.Fatalf("expected *httpx.AppError, got %T: %v", err, err)
	}
	if ae.Status != 400 {
		t.Errorf("expected status 400, got %d", ae.Status)
	}
	if ae.Code != "bad_request" {
		t.Errorf("expected code %q, got %q", "bad_request", ae.Code)
	}
}
