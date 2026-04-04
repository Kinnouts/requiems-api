package swift

import (
	"net/http"
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
	// Location code with digit is valid.
	got, err := sanitizeSWIFT("DEUT2EDB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEUT2EDBXXX" {
		t.Errorf("expected DEUT2EDBXXX, got %q", got)
	}
}

func TestSanitizeSWIFT_TooShort(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDE")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_TooLong(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB001X")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_9Chars(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB0")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_10Chars(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB01")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_Empty(t *testing.T) {
	_, err := sanitizeSWIFT("")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_DigitInBankCode(t *testing.T) {
	_, err := sanitizeSWIFT("1EUTDEDB")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_DigitInCountryCode(t *testing.T) {
	_, err := sanitizeSWIFT("DEUT1EDB")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
}

func TestSanitizeSWIFT_InvalidBranchCode(t *testing.T) {
	_, err := sanitizeSWIFT("DEUTDEDB0-1")
	assertAppError(t, err, http.StatusBadRequest, "bad_request")
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

// assertAppError verifies that err is an *httpx.AppError with the expected
// HTTP status and error code.
func assertAppError(t *testing.T, err error, status int, code string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	ae, ok := err.(*httpx.AppError)
	if !ok {
		t.Fatalf("expected *httpx.AppError, got %T: %v", err, err)
	}
	if ae.Status != status {
		t.Errorf("expected status %d, got %d", status, ae.Status)
	}
	if ae.Code != code {
		t.Errorf("expected code %q, got %q", code, ae.Code)
	}
}
