package detectlanguage

import (
	"testing"
)

func TestService_Detect_French(t *testing.T) {
	svc := NewService()

	result := svc.Detect("Bonjour, comment ça va?")

	if result.Language != "French" {
		t.Errorf("expected language to be French, got %q", result.Language)
	}
	if result.Code != "fr" {
		t.Errorf("expected code to be fr, got %q", result.Code)
	}
	if result.Confidence <= 0 {
		t.Error("expected confidence to be greater than 0")
	}
}

func TestService_Detect_English(t *testing.T) {
	svc := NewService()

	result := svc.Detect("The quick brown fox jumps over the lazy dog")

	if result.Language != "English" {
		t.Errorf("expected language to be English, got %q", result.Language)
	}
	if result.Code != "en" {
		t.Errorf("expected code to be en, got %q", result.Code)
	}
	if result.Confidence <= 0 {
		t.Error("expected confidence to be greater than 0")
	}
}

func TestService_Detect_Spanish(t *testing.T) {
	svc := NewService()

	result := svc.Detect("El rápido zorro marrón salta sobre el perro perezoso")

	if result.Language != "Spanish" {
		t.Errorf("expected language to be Spanish, got %q", result.Language)
	}
	if result.Code != "es" {
		t.Errorf("expected code to be es, got %q", result.Code)
	}
	if result.Confidence <= 0 {
		t.Error("expected confidence to be greater than 0")
	}
}

func TestService_Detect_ConfidenceRange(t *testing.T) {
	svc := NewService()

	result := svc.Detect("This is a longer English sentence to ensure reliable detection.")

	if result.Confidence < 0 || result.Confidence > 1 {
		t.Errorf("confidence %f is outside expected range [0, 1]", result.Confidence)
	}
}
