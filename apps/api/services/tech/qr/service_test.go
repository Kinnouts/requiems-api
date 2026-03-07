package qr

import (
	"testing"
)

func TestService_Generate_ValidData(t *testing.T) {
	svc := NewService()

	png, err := svc.Generate("https://example.com", 256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes")
	}

	// Verify PNG signature (\x89PNG\r\n\x1a\n)
	if len(png) < 8 || string(png[:4]) != "\x89PNG" {
		t.Error("expected valid PNG signature")
	}
}

func TestService_Generate_SmallSize(t *testing.T) {
	svc := NewService()

	png, err := svc.Generate("hello", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes for small size")
	}
}

func TestService_Generate_LargeSize(t *testing.T) {
	svc := NewService()

	png, err := svc.Generate("https://example.com/very/long/path?foo=bar&baz=qux", 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes for large size")
	}
}

func TestService_Generate_EmptyData(t *testing.T) {
	svc := NewService()

	_, err := svc.Generate("", 256)
	if err == nil {
		t.Error("expected error for empty data")
	}
}
