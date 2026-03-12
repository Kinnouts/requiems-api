package bin

import (
	"testing"
)

// ---- sanitizeBIN ----

func TestSanitizeBIN_Valid6Digit(t *testing.T) {
	got, err := sanitizeBIN("424242")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "424242" {
		t.Errorf("got %q, want %q", got, "424242")
	}
}

func TestSanitizeBIN_Valid8Digit(t *testing.T) {
	got, err := sanitizeBIN("42424242")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "42424242" {
		t.Errorf("got %q, want %q", got, "42424242")
	}
}

func TestSanitizeBIN_StripsDashes(t *testing.T) {
	got, err := sanitizeBIN("4242-42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "424242" {
		t.Errorf("got %q, want %q", got, "424242")
	}
}

func TestSanitizeBIN_StripsSpaces(t *testing.T) {
	got, err := sanitizeBIN("  4242 42  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "424242" {
		t.Errorf("got %q, want %q", got, "424242")
	}
}

func TestSanitizeBIN_TooShort(t *testing.T) {
	_, err := sanitizeBIN("42424")
	if err == nil {
		t.Fatal("expected error for 5-digit BIN, got nil")
	}
}

func TestSanitizeBIN_TooLong(t *testing.T) {
	_, err := sanitizeBIN("424242424")
	if err == nil {
		t.Fatal("expected error for 9-digit BIN, got nil")
	}
}

func TestSanitizeBIN_Empty(t *testing.T) {
	_, err := sanitizeBIN("")
	if err == nil {
		t.Fatal("expected error for empty BIN, got nil")
	}
}

func TestSanitizeBIN_NonDigits(t *testing.T) {
	_, err := sanitizeBIN("abcdef")
	if err == nil {
		t.Fatal("expected error for non-digit BIN, got nil")
	}
}

func TestSanitizeBIN_MixedAlphaDigits(t *testing.T) {
	_, err := sanitizeBIN("4242ab")
	if err == nil {
		t.Fatal("expected error for mixed alpha-digit BIN, got nil")
	}
}

// ---- luhnValid ----

func TestLuhnValid_KnownValid(t *testing.T) {
	// 424242 — well-known Visa test BIN with valid Luhn
	if !luhnValid("424242") {
		t.Error("expected luhnValid(424242) = true")
	}
}

func TestLuhnValid_KnownInvalid(t *testing.T) {
	if luhnValid("123456") {
		t.Error("expected luhnValid(123456) = false")
	}
}

func TestLuhnValid_AllZeros(t *testing.T) {
	// 000000 → sum = 0, 0 % 10 = 0 → valid
	if !luhnValid("000000") {
		t.Error("expected luhnValid(000000) = true")
	}
}

func TestLuhnValid_8DigitValid(t *testing.T) {
	if !luhnValid("42424242") {
		t.Error("expected luhnValid(42424242) = true")
	}
}

// ---- detectScheme ----

func TestDetectScheme_Visa(t *testing.T) {
	cases := []string{"424242", "400000", "499999"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "visa" {
			t.Errorf("detectScheme(%s) = %q, want %q", bin, got, "visa")
		}
	}
}

func TestDetectScheme_Mastercard5Series(t *testing.T) {
	cases := map[string]string{
		"510000": "mastercard",
		"520000": "mastercard",
		"530000": "mastercard",
		"540000": "mastercard",
		"550000": "mastercard",
	}
	for bin, want := range cases {
		if got := detectScheme(bin); got != want {
			t.Errorf("detectScheme(%s) = %q, want %q", bin, got, want)
		}
	}
}

func TestDetectScheme_Mastercard2Series(t *testing.T) {
	cases := []string{"222100", "272000", "250000"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "mastercard" {
			t.Errorf("detectScheme(%s) = %q, want %q", bin, got, "mastercard")
		}
	}
}

func TestDetectScheme_Mastercard2SeriesBoundaryLow(t *testing.T) {
	// 2220xx is NOT Mastercard (range starts at 2221)
	if got := detectScheme("222099"); got == "mastercard" {
		t.Errorf("detectScheme(222099) should not be mastercard, got %q", got)
	}
}

func TestDetectScheme_Mastercard2SeriesBoundaryHigh(t *testing.T) {
	// 2721xx is NOT Mastercard (range ends at 2720)
	if got := detectScheme("272100"); got == "mastercard" {
		t.Errorf("detectScheme(272100) should not be mastercard, got %q", got)
	}
}

func TestDetectScheme_Amex(t *testing.T) {
	cases := []string{"340000", "370000", "378282"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "amex" {
			t.Errorf("detectScheme(%s) = %q, want %q", bin, got, "amex")
		}
	}
}

func TestDetectScheme_Discover6011(t *testing.T) {
	if got := detectScheme("601100"); got != "discover" {
		t.Errorf("detectScheme(601100) = %q, want discover", got)
	}
}

func TestDetectScheme_Discover622Inside(t *testing.T) {
	cases := []string{"622126", "622500", "622925"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "discover" {
			t.Errorf("detectScheme(%s) = %q, want discover", bin, got)
		}
	}
}

func TestDetectScheme_Discover622BoundaryLow(t *testing.T) {
	// 622125 is below the Discover range → UnionPay
	if got := detectScheme("622125"); got == "discover" {
		t.Errorf("detectScheme(622125) should not be discover, got %q", got)
	}
}

func TestDetectScheme_Discover622BoundaryHigh(t *testing.T) {
	// 622926 is above the Discover range → UnionPay
	if got := detectScheme("622926"); got == "discover" {
		t.Errorf("detectScheme(622926) should not be discover, got %q", got)
	}
}

func TestDetectScheme_Discover65(t *testing.T) {
	if got := detectScheme("650000"); got != "discover" {
		t.Errorf("detectScheme(650000) = %q, want discover", got)
	}
}

func TestDetectScheme_JCB(t *testing.T) {
	cases := []string{"352800", "358900", "356000"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "jcb" {
			t.Errorf("detectScheme(%s) = %q, want jcb", bin, got)
		}
	}
}

func TestDetectScheme_Diners(t *testing.T) {
	cases := []string{"300000", "305999", "360000", "380000"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "diners" {
			t.Errorf("detectScheme(%s) = %q, want diners", bin, got)
		}
	}
}

func TestDetectScheme_UnionPay(t *testing.T) {
	cases := []string{"620000", "810000"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "unionpay" {
			t.Errorf("detectScheme(%s) = %q, want unionpay", bin, got)
		}
	}
}

func TestDetectScheme_Maestro(t *testing.T) {
	cases := []string{"630400", "675900", "676100", "676200", "676300"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "maestro" {
			t.Errorf("detectScheme(%s) = %q, want maestro", bin, got)
		}
	}
}

func TestDetectScheme_Mir(t *testing.T) {
	cases := []string{"220000", "220100", "220200", "220300", "220400"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "mir" {
			t.Errorf("detectScheme(%s) = %q, want mir", bin, got)
		}
	}
}

func TestDetectScheme_MirVsMastercardBoundary(t *testing.T) {
	// 2205xx is between Mir (≤2204) and Mastercard 2-series (≥2221): neither
	got := detectScheme("220500")
	if got == "mir" || got == "mastercard" {
		t.Errorf("detectScheme(220500) = %q, expected neither mir nor mastercard", got)
	}
}

func TestDetectScheme_RuPay(t *testing.T) {
	cases := []string{"600000", "652100", "652200"}
	for _, bin := range cases {
		if got := detectScheme(bin); got != "rupay" {
			t.Errorf("detectScheme(%s) = %q, want rupay", bin, got)
		}
	}
}

func TestDetectScheme_Unknown(t *testing.T) {
	if got := detectScheme("999999"); got != "" {
		t.Errorf("detectScheme(999999) = %q, want empty string", got)
	}
}

func TestDetectScheme_TooShort(t *testing.T) {
	if got := detectScheme("42"); got != "" {
		t.Errorf("detectScheme(42) = %q, want empty string for too-short input", got)
	}
}
