package holidays

import (
	"testing"
)

func TestService_GetHolidays_US_2025(t *testing.T) {
	svc := NewService()

	resp, err := svc.GetHolidays("US", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Country != "US" {
		t.Errorf("expected country 'US', got %q", resp.Country)
	}
	if resp.Year != 2025 {
		t.Errorf("expected year 2025, got %d", resp.Year)
	}
	if len(resp.Holidays) == 0 {
		t.Error("expected non-empty holidays list")
	}
}

func TestService_GetHolidays_GB_2025(t *testing.T) {
	svc := NewService()

	resp, err := svc.GetHolidays("GB", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Country != "GB" {
		t.Errorf("expected country 'GB', got %q", resp.Country)
	}
	if len(resp.Holidays) == 0 {
		t.Error("expected non-empty holidays list")
	}
}

func TestService_GetHolidays_Japan(t *testing.T) {
	svc := NewService()

	resp, err := svc.GetHolidays("JP", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Holidays) == 0 {
		t.Error("expected non-empty holidays list for Japan")
	}
}

func TestService_GetHolidays_Germany(t *testing.T) {
	svc := NewService()

	resp, err := svc.GetHolidays("DE", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Holidays) == 0 {
		t.Error("expected non-empty holidays list for Germany")
	}
}

func TestService_GetHolidays_NewYear(t *testing.T) {
	svc := NewService()

	resp, err := svc.GetHolidays("US", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, h := range resp.Holidays {
		if h.Name == "New Year's Day" && h.Date == "2025-01-01" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected New Year's Day 2025-01-01 in US holidays")
	}
}

func TestService_GetHolidays_DateFormat(t *testing.T) {
	svc := NewService()

	resp, err := svc.GetHolidays("US", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, h := range resp.Holidays {
		if len(h.Date) != 10 {
			t.Errorf("expected date in YYYY-MM-DD format, got %q", h.Date)
		}
		if h.Date[4] != '-' || h.Date[7] != '-' {
			t.Errorf("expected date in YYYY-MM-DD format, got %q", h.Date)
		}
	}
}

func TestService_GetHolidays_InvalidCountry(t *testing.T) {
	svc := NewService()

	_, err := svc.GetHolidays("XX", 2025)
	if err == nil {
		t.Fatal("expected error for invalid country code, got nil")
	}
	if err != nil {
		if !contains(err.Error(), "no holidays found") {
			t.Errorf("expected error message to contain 'no holidays found', got %q", err.Error())
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
