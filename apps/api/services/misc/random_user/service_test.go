package randomuser

import (
	"strings"
	"testing"
)

func TestGenerate_FieldsPopulated(t *testing.T) {
	svc := NewService()

	for range 20 {
		u := svc.Generate()

		if u.Name == "" {
			t.Error("Name must not be empty")
		}
		if !strings.Contains(u.Name, " ") {
			t.Errorf("Name should be first + last: got %q", u.Name)
		}
		if u.Email == "" {
			t.Error("Email must not be empty")
		}
		if !strings.Contains(u.Email, "@") {
			t.Errorf("Email missing @: got %q", u.Email)
		}
		if u.Phone == "" {
			t.Error("Phone must not be empty")
		}
		if !strings.HasPrefix(u.Phone, "+1-555-") {
			t.Errorf("Phone should start with +1-555-: got %q", u.Phone)
		}
		if u.Address.Street == "" {
			t.Error("Address.Street must not be empty")
		}
		if u.Address.City == "" {
			t.Error("Address.City must not be empty")
		}
		if u.Address.State == "" {
			t.Error("Address.State must not be empty")
		}
		if len(u.Address.Zip) != 5 {
			t.Errorf("Zip should be 5 digits: got %q", u.Address.Zip)
		}
		if u.Address.Country != "United States" {
			t.Errorf("Country should be United States: got %q", u.Address.Country)
		}
		if !strings.HasPrefix(u.Avatar, "https://api.dicebear.com/") {
			t.Errorf("Avatar should start with dicebear URL: got %q", u.Avatar)
		}
	}
}

func TestGenerate_EmailMatchesName(t *testing.T) {
	svc := NewService()

	for range 20 {
		u := svc.Generate()

		parts := strings.SplitN(u.Name, " ", 2)
		if len(parts) != 2 {
			t.Fatalf("expected two name parts, got: %q", u.Name)
		}

		expectedPrefix := strings.ToLower(parts[0] + "." + parts[1]) + "@"
		if !strings.HasPrefix(u.Email, expectedPrefix) {
			t.Errorf("email %q should start with %q", u.Email, expectedPrefix)
		}
	}
}

func TestGenerate_AvatarContainsName(t *testing.T) {
	svc := NewService()

	u := svc.Generate()
	if !strings.Contains(u.Avatar, "seed=") {
		t.Errorf("avatar URL should contain seed parameter: %q", u.Avatar)
	}
}
