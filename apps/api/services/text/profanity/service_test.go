package profanity

import (
	"context"
	"testing"
)

func TestService_Check_NoProfanity(t *testing.T) {
	svc := NewService()

	result := svc.Check(context.Background(),"Hello, world!")

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

	result := svc.Check(context.Background(),"What the fuck is this shit")

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

	// go-away detects "shit" as a substring of "BULLSHIT" and censors only
	// the matched portion; the flagged canonical word is "shit".
	result := svc.Check(context.Background(),"This is BULLSHIT")

	if !result.HasProfanity {
		t.Error("expected HasProfanity to be true for uppercase word")
	}
	if len(result.FlaggedWords) != 1 || result.FlaggedWords[0] != "shit" {
		t.Errorf("expected flagged word [\"shit\"], got %v", result.FlaggedWords)
	}
}

func TestService_Check_DeduplicatesFlaggedWords(t *testing.T) {
	svc := NewService()

	result := svc.Check(context.Background(),"shit shit shit")

	if len(result.FlaggedWords) != 1 {
		t.Errorf("expected 1 unique flagged word, got %d: %v", len(result.FlaggedWords), result.FlaggedWords)
	}
}

func TestService_Check_EmptyFlaggedWordsSlice(t *testing.T) {
	svc := NewService()

	result := svc.Check(context.Background(),"clean text here")

	// FlaggedWords must be an empty slice, not nil (for consistent JSON serialisation).
	if result.FlaggedWords == nil {
		t.Error("expected FlaggedWords to be an empty slice, not nil")
	}
}

func TestService_Check_EmptyText(t *testing.T) {
	svc := NewService()

	result := svc.Check(context.Background(),"")

	if result.HasProfanity {
		t.Error("expected no profanity for empty text")
	}
	if result.Censored != "" {
		t.Errorf("expected empty censored, got %q", result.Censored)
	}
}

func TestService_Check_LeetSpeak(t *testing.T) {
	svc := NewService()

	// go-away handles leet-speak obfuscation out of the box.
	result := svc.Check(context.Background(),"F   u   C  k th1$ $h!t")

	if !result.HasProfanity {
		t.Error("expected HasProfanity to be true for leet-speak input")
	}
	if len(result.FlaggedWords) == 0 {
		t.Error("expected at least one flagged word for leet-speak input")
	}
}
