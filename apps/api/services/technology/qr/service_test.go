package qr

import (
	"testing"
)

func TestService_Generate_ValidData(t *testing.T) {
	svc := NewService()

	png, err := svc.Generate("https://example.com", 256, "medium")
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

	png, err := svc.Generate("hello", 50, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes for small size")
	}
}

func TestService_Generate_LargeSize(t *testing.T) {
	svc := NewService()

	png, err := svc.Generate("https://example.com/very/long/path?foo=bar&baz=qux", 1000, "low")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes for large size")
	}
}

func TestService_Generate_HighestRecovery(t *testing.T) {
	svc := NewService()

	png, err := svc.Generate("https://example.com", 256, "highest")
	if err != nil {
		t.Fatalf("unexpected error for highest recovery level: %v", err)
	}

	if len(png) == 0 {
		t.Error("expected non-empty PNG bytes for highest recovery level")
	}
}

func TestService_Generate_AllRecoveryLevels(t *testing.T) {
	svc := NewService()

	levels := []string{"low", "medium", "high", "highest"}
	for _, level := range levels {
		png, err := svc.Generate("test", 256, level)
		if err != nil {
			t.Errorf("unexpected error for recovery=%q: %v", level, err)
			continue
		}
		if len(png) == 0 {
			t.Errorf("expected non-empty PNG for recovery=%q", level)
		}
	}
}

func TestService_Generate_EmptyData(t *testing.T) {
	svc := NewService()

	_, err := svc.Generate("", 256, "medium")
	if err == nil {
		t.Error("expected error for empty data")
	}
}
