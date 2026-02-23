package profanity

import (
	"testing"
)

func TestService_Check_NoProfanity(t *testing.T) {
	svc := NewService()

	result := svc.Check("Hello, world!")

	if result.HasProfanity {
		t.Error("expected HasProfanity to be false")
	}
	if result.Censored != "Hello, world!" {
		t.Errorf("expected censored to equal input, got %q", result.Censored)
	}
	if len(result.FlaggedWords) != 0 {
		t.Errorf("expected no flagged words, got %v", result.FlaggedWords)
	}
}

func TestService_Check_WithProfanity(t *testing.T) {
	svc := NewService()

	result := svc.Check("What the fuck is this shit")

	if !result.HasProfanity {
		t.Error("expected HasProfanity to be true")
	}
	if result.Censored != "What the **** is this ****" {
		t.Errorf("unexpected censored output: %q", result.Censored)
	}
	if len(result.FlaggedWords) != 2 {
		t.Errorf("expected 2 flagged words, got %d: %v", len(result.FlaggedWords), result.FlaggedWords)
	}
}

func TestService_Check_CaseInsensitive(t *testing.T) {
	svc := NewService()

	result := svc.Check("This is BULLSHIT")

	if !result.HasProfanity {
		t.Error("expected HasProfanity to be true for uppercase word")
	}
	if result.Censored != "This is ********" {
		t.Errorf("unexpected censored output: %q", result.Censored)
	}
	if len(result.FlaggedWords) != 1 || result.FlaggedWords[0] != "bullshit" {
		t.Errorf("expected flagged word 'bullshit', got %v", result.FlaggedWords)
	}
}

func TestService_Check_DeduplicatesFlaggedWords(t *testing.T) {
	svc := NewService()

	result := svc.Check("shit shit shit")

	if len(result.FlaggedWords) != 1 {
		t.Errorf("expected 1 unique flagged word, got %d: %v", len(result.FlaggedWords), result.FlaggedWords)
	}
}

func TestService_Check_EmptyFlaggedWordsSlice(t *testing.T) {
	svc := NewService()

	result := svc.Check("clean text here")

	// FlaggedWords must be an empty slice, not nil (for consistent JSON serialisation).
	if result.FlaggedWords == nil {
		t.Error("expected FlaggedWords to be an empty slice, not nil")
	}
}

func TestService_Check_EmptyText(t *testing.T) {
	svc := NewService()

	result := svc.Check("")

	if result.HasProfanity {
		t.Error("expected no profanity for empty text")
	}
	if result.Censored != "" {
		t.Errorf("expected empty censored, got %q", result.Censored)
	}
}

func TestService_Check_PunctuationAroundProfanity(t *testing.T) {
	svc := NewService()

	// Punctuation should not prevent detection.
	result := svc.Check("Oh, damn!")

	if !result.HasProfanity {
		t.Error("expected HasProfanity to be true")
	}
	if result.Censored != "Oh, ****!" {
		t.Errorf("unexpected censored output: %q", result.Censored)
	}
}
