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
