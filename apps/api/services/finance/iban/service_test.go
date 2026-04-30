package iban

import "testing"

// ---- normalizeIBAN ----

func TestNormalizeIBAN_StripsSpaces(t *testing.T) {
	got := normalizeIBAN("DE89 3704 0044 0532 0130 00")
	want := "DE89370400440532013000"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeIBAN_Uppercases(t *testing.T) {
	got := normalizeIBAN("de89370400440532013000")
	want := "DE89370400440532013000"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeIBAN_TrimsTrimSpace(t *testing.T) {
	got := normalizeIBAN("  DE89370400440532013000  ")
	want := "DE89370400440532013000"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// ---- basicFormatOK ----

func TestBasicFormatOK_ValidDE(t *testing.T) {
	if !basicFormatOK("DE89370400440532013000") {
		t.Error("expected true for valid German IBAN")
	}
}

func TestBasicFormatOK_ValidGB(t *testing.T) {
	if !basicFormatOK("GB82WEST12345698765432") {
		t.Error("expected true for valid UK IBAN")
	}
}

func TestBasicFormatOK_TooShort(t *testing.T) {
	if basicFormatOK("DE89") {
		t.Error("expected false for 4-char input")
	}
}

func TestBasicFormatOK_Empty(t *testing.T) {
	if basicFormatOK("") {
		t.Error("expected false for empty string")
	}
}

func TestBasicFormatOK_DigitInCountryCode(t *testing.T) {
	if basicFormatOK("1E89370400440532013000") {
		t.Error("expected false when first char is a digit")
	}
}

func TestBasicFormatOK_LetterInCheckDigits(t *testing.T) {
	if basicFormatOK("DEAB370400440532013000") {
		t.Error("expected false when check digits contain letters")
	}
}

func TestBasicFormatOK_SpecialCharacter(t *testing.T) {
	if basicFormatOK("DE89!70400440532013000") {
		t.Error("expected false for input containing '!'")
	}
}

// ---- validateChecksum ----

var validIBANs = []struct {
	iban    string
	country string
}{
	{"DE89370400440532013000", "Germany"},
	{"GB82WEST12345698765432", "United Kingdom"},
	{"FR7630006000011234567890189", "France"},
	{"NL91ABNA0417164300", "Netherlands"},
	{"CH9300762011623852957", "Switzerland"},
	{"AT611904300234573201", "Austria"},
	{"BE68539007547034", "Belgium"},
	{"PL61109010140000071219812874", "Poland"},
}

func TestValidateChecksum_KnownValidIBANs(t *testing.T) {
	for _, tc := range validIBANs {
		if !validateChecksum(tc.iban) {
			t.Errorf("validateChecksum(%s) = false, expected true (%s)", tc.iban, tc.country)
		}
	}
}

func TestValidateChecksum_WrongCheckDigits(t *testing.T) {
	// DE89... with check digits changed to 00 — should fail.
	if validateChecksum("DE00370400440532013000") {
		t.Error("expected false for IBAN with wrong check digits (DE00...)")
	}
}

func TestValidateChecksum_TransposedDigits(t *testing.T) {
	// Swapping two adjacent digits in the BBAN breaks the checksum.
	if validateChecksum("DE89370400440352013000") {
		t.Error("expected false for IBAN with transposed digits in BBAN")
	}
}

func TestValidateChecksum_AllZeroCheckDigits(t *testing.T) {
	if validateChecksum("GB00WEST12345698765432") {
		t.Error("expected false for 00 check digits")
	}
}

// ---- mod97 ----

func TestMod97_DEExample(t *testing.T) {
	// DE89370400440532013000 rearranged + letters replaced:
	// BBAN+CC+CD → 370400440532013000DE89 → 370400440532013000131489
	if got := mod97("370400440532013000131489"); got != 1 {
		t.Errorf("mod97 = %d, want 1", got)
	}
}

func TestMod97_NLExample(t *testing.T) {
	// NL91ABNA0417164300 rearranged → ABNA0417164300NL91
	// A=10 B=11 N=23 A=10 | 0417164300 | N=23 L=21 | 91
	// → "101123100417164300232191"
	if got := mod97("101123100417164300232191"); got != 1 {
		t.Errorf("mod97 = %d, want 1", got)
	}
}

// ---- extract ----

func TestExtract_HappyPath(t *testing.T) {
	got := extract("ABCDEFGH", 2, 3)
	if got != "CDE" {
		t.Errorf("extract = %q, want %q", got, "CDE")
	}
}

func TestExtract_FromStart(t *testing.T) {
	got := extract("37040044ABCDE", 0, 8)
	if got != "37040044" {
		t.Errorf("extract = %q, want %q", got, "37040044")
	}
}

func TestExtract_ZeroLength(t *testing.T) {
	if got := extract("ABCDEF", 0, 0); got != "" {
		t.Errorf("extract with 0 length should be empty, got %q", got)
	}
}

func TestExtract_OutOfBounds(t *testing.T) {
	if got := extract("ABC", 1, 10); got != "" {
		t.Errorf("out-of-bounds extract should be empty, got %q", got)
	}
}

func TestExtract_NegativeOffset(t *testing.T) {
	if got := extract("ABC", -1, 2); got != "" {
		t.Errorf("negative offset extract should be empty, got %q", got)
	}
}
