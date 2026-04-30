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
		if u.Email == "" {
			t.Error("Email must not be empty")
		}
		if !strings.Contains(u.Email, "@") {
			t.Errorf("Email missing @: got %q", u.Email)
		}
		if u.Phone == "" {
			t.Error("Phone must not be empty")
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
		if u.Address.Zip == "" {
			t.Error("Address.Zip must not be empty")
		}
		if u.Address.Country == "" {
			t.Error("Address.Country must not be empty")
		}
		if !strings.HasPrefix(u.Avatar, "https://api.dicebear.com/") {
			t.Errorf("Avatar should start with dicebear URL: got %q", u.Avatar)
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
