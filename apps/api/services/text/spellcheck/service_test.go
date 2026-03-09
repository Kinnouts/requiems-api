package spellcheck

import (
	"testing"
)

func TestService_Check_NoMistakes(t *testing.T) {
	svc := NewService()

	result := svc.Check("This is a test")

	if result.Corrected != "This is a test" {
		t.Errorf("expected corrected to equal input, got %q", result.Corrected)
	}
	if len(result.Corrections) != 0 {
		t.Errorf("expected no corrections, got %v", result.Corrections)
	}
}

func TestService_Check_MisspelledWords(t *testing.T) {
	svc := NewService()

	result := svc.Check("Ths is a tset")

	if len(result.Corrections) == 0 {
		t.Fatal("expected corrections, got none")
	}

	foundThs := false
	foundTset := false
	for _, c := range result.Corrections {
		if c.Original == "Ths" && c.Position == 0 && c.Suggested != "" {
			foundThs = true
		}
		if c.Original == "tset" && c.Position == 9 && c.Suggested != "" {
			foundTset = true
		}
	}

	if !foundThs {
		t.Errorf("expected correction for Ths at position 0; got %+v", result.Corrections)
	}
	if !foundTset {
		t.Errorf("expected correction for tset at position 9; got %+v", result.Corrections)
	}
}

func TestService_Check_EmptyText(t *testing.T) {
	svc := NewService()

	// Empty text is not a valid request (validate:"required" enforces that at
	// the HTTP layer), but the service itself should return a safe empty result.
	result := svc.Check("")

	if result.Corrected != "" {
		t.Errorf("expected empty corrected string, got %q", result.Corrected)
	}
	if len(result.Corrections) != 0 {
		t.Errorf("expected no corrections for empty input, got %v", result.Corrections)
	}
}

func TestService_Check_CorrectedTextReflectsFixes(t *testing.T) {
	svc := NewService()

	result := svc.Check("Ths is a tset")

	if result.Corrected == "Ths is a tset" {
		t.Error("expected corrected text to differ from misspelled input")
	}
}

func TestService_Check_CorrectionsSliceNotNil(t *testing.T) {
	svc := NewService()

	result := svc.Check("Hello world")

	if result.Corrections == nil {
		t.Error("expected non-nil corrections slice for clean input")
	}
}

func TestMatchCase_LowerInput(t *testing.T) {
	got := matchCase("abc", "suggested")
	if got != "suggested" {
		t.Errorf("expected %q, got %q", "suggested", got)
	}
}

func TestMatchCase_CapitalisedInput(t *testing.T) {
	got := matchCase("Abc", "suggested")
	if got != "Suggested" {
		t.Errorf("expected %q, got %q", "Suggested", got)
	}
}

func TestMatchCase_AllUpperInput(t *testing.T) {
	got := matchCase("ABC", "suggested")
	if got != "SUGGESTED" {
		t.Errorf("expected %q, got %q", "SUGGESTED", got)
	}
}
