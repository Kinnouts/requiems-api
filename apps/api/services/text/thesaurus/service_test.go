package thesaurus

import (
	"testing"
)

func TestService_Lookup_KnownWord(t *testing.T) {
	svc := NewService()

	result, err := svc.Lookup("happy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Word != "happy" {
		t.Errorf("expected word %q, got %q", "happy", result.Word)
	}

	if len(result.Synonyms) == 0 {
		t.Error("expected at least one synonym")
	}

	if len(result.Antonyms) == 0 {
		t.Error("expected at least one antonym")
	}
}

func TestService_Lookup_CaseInsensitive(t *testing.T) {
	svc := NewService()

	tests := []string{"HAPPY", "Happy", "hApPy", "happy"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result, err := svc.Lookup(input)
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", input, err)
			}
			if result.Word != "happy" {
				t.Errorf("expected word %q, got %q", "happy", result.Word)
			}
		})
	}
}

func TestService_Lookup_UnknownWord(t *testing.T) {
	svc := NewService()

	_, err := svc.Lookup("zzyzx")
	if err == nil {
		t.Fatal("expected error for unknown word, got nil")
	}
}

func TestService_Lookup_ReturnsNonNilSlices(t *testing.T) {
	svc := NewService()

	result, err := svc.Lookup("happy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Synonyms == nil {
		t.Error("expected non-nil synonyms slice")
	}

	if result.Antonyms == nil {
		t.Error("expected non-nil antonyms slice")
	}
}
