package normalize

import (
	"testing"

	normalizer "github.com/bobadilla-tech/go-email-normalizer"
)

// containsChange reports whether target appears in the changes slice.
func containsChange(changes []normalizer.Change, target normalizer.Change) bool {
	for _, c := range changes {
		if c == target {
			return true
		}
	}
	return false
}

func TestService_Normalize_OriginalPreserved(t *testing.T) {
	svc := NewService()

	input := "User@Example.com"
	result := svc.Normalize(input)

	if result.Original != input {
		t.Errorf("expected Original %q, got %q", input, result.Original)
	}
}

func TestService_Normalize_SplitsLocalAndDomain(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("user@example.com")

	if result.Local != "user" {
		t.Errorf("expected Local %q, got %q", "user", result.Local)
	}
	if result.Domain != "example.com" {
		t.Errorf("expected Domain %q, got %q", "example.com", result.Domain)
	}
}

func TestService_Normalize_LowercasesDomain(t *testing.T) {
	svc := NewService()

	// For unknown providers the local part is preserved (case-sensitive per
	// RFC 5321); only the domain is lowercased.
	result := svc.Normalize("User@Example.com")

	if result.Normalized != "User@example.com" {
		t.Errorf("expected normalized %q, got %q", "User@example.com", result.Normalized)
	}
	if !containsChange(result.Changes, normalizer.ChangeLowercase) {
		t.Errorf("expected ChangeLowercase in changes, got %v", result.Changes)
	}
}

func TestService_Normalize_NoChangesForAlreadyNormalisedEmail(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("user@example.com")

	if result.Normalized != "user@example.com" {
		t.Errorf("expected normalized to equal input, got %q", result.Normalized)
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected no changes for already-normalised email, got %v", result.Changes)
	}
}

func TestService_Normalize_GmailRemovesDots(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("te.st.user@gmail.com")

	if result.Normalized != "testuser@gmail.com" {
		t.Errorf("expected normalized %q, got %q", "testuser@gmail.com", result.Normalized)
	}
	if !containsChange(result.Changes, normalizer.ChangeRemovedDots) {
		t.Errorf("expected ChangeRemovedDots in changes, got %v", result.Changes)
	}
}

func TestService_Normalize_GmailRemovesPlusTag(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("testuser+spam@gmail.com")

	if result.Normalized != "testuser@gmail.com" {
		t.Errorf("expected normalized %q, got %q", "testuser@gmail.com", result.Normalized)
	}
	if !containsChange(result.Changes, normalizer.ChangeRemovedPlusTag) {
		t.Errorf("expected ChangeRemovedPlusTag in changes, got %v", result.Changes)
	}
}

func TestService_Normalize_GooglemailCanonicalisedToGmail(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("user@googlemail.com")

	if result.Normalized != "user@gmail.com" {
		t.Errorf("expected normalized %q, got %q", "user@gmail.com", result.Normalized)
	}
	if !containsChange(result.Changes, normalizer.ChangeCanonicalisedDomain) {
		t.Errorf("expected ChangeCanonicalisedDomain in changes, got %v", result.Changes)
	}
}

func TestService_Normalize_WhitespaceTrimmed(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("  user@example.com  ")

	if result.Normalized != "user@example.com" {
		t.Errorf("expected whitespace to be trimmed, got %q", result.Normalized)
	}
	if !containsChange(result.Changes, normalizer.ChangeTrimmedWhitespace) {
		t.Errorf("expected ChangeTrimmedWhitespace in changes, got %v", result.Changes)
	}
}

func TestService_Normalize_NormalizedFieldPopulated(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("user@example.com")

	if result.Normalized == "" {
		t.Error("expected non-empty Normalized field")
	}
}

func TestService_Normalize_LocalAndDomainMatchNormalized(t *testing.T) {
	svc := NewService()

	result := svc.Normalize("Test.User+tag@Gmail.com")

	if result.Local+"@"+result.Domain != result.Normalized {
		t.Errorf("expected Local@Domain to equal Normalized %q, got %q@%q",
			result.Normalized, result.Local, result.Domain)
	}
}
