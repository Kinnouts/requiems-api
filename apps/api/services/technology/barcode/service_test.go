package barcode

import (
	"testing"
)

func TestService_Generate_Code128(t *testing.T) {
	svc := NewService()

	png, width, height, err := svc.Generate("123456789", "code128")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes")
	}

	if string(png[:4]) != "\x89PNG" {
		t.Error("expected valid PNG signature")
	}

	if width != defaultWidth || height != defaultHeight {
		t.Errorf("expected %dx%d, got %dx%d", defaultWidth, defaultHeight, width, height)
	}
}

func TestService_Generate_Code93(t *testing.T) {
	svc := NewService()

	png, _, _, err := svc.Generate("HELLO", "code93")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes")
	}
}

func TestService_Generate_Code39(t *testing.T) {
	svc := NewService()

	png, _, _, err := svc.Generate("HELLO123", "code39")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes")
	}
}

func TestService_Generate_EAN8(t *testing.T) {
	svc := NewService()

	// 7 digits (checksum auto-calculated)
	png, _, _, err := svc.Generate("1234567", "ean8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes")
	}
}

func TestService_Generate_EAN13(t *testing.T) {
	svc := NewService()

	// 12 digits (checksum auto-calculated)
	png, _, _, err := svc.Generate("123456789012", "ean13")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes")
	}
}

func TestService_Generate_EAN8_InvalidLength(t *testing.T) {
	svc := NewService()

	_, _, _, err := svc.Generate("123", "ean8")
	if err == nil {
		t.Error("expected error for invalid EAN-8 length")
	}
}

func TestService_Generate_EAN13_InvalidLength(t *testing.T) {
	svc := NewService()

	_, _, _, err := svc.Generate("123456789", "ean13")
	if err == nil {
		t.Error("expected error for invalid EAN-13 length")
	}
}

func TestService_Generate_UnsupportedType(t *testing.T) {
	svc := NewService()

	_, _, _, err := svc.Generate("hello", "qrcode")
	if err == nil {
		t.Error("expected error for unsupported barcode type")
	}
}

func TestService_Generate_EmptyData_Code128(t *testing.T) {
	svc := NewService()

	_, _, _, err := svc.Generate("", "code128")
	if err == nil {
		t.Error("expected error for empty data")
	}
}
