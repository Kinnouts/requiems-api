package holidays

import (
	"strings"
	"testing"
)

func TestService_GetHolidays_Countries(t *testing.T) {
	cases := []struct {
		name    string
		country string
		year    int
	}{
		{name: "US_2025", country: "US", year: 2025},
		{name: "GB_2025", country: "GB", year: 2025},
		{name: "Japan", country: "JP", year: 2025},
		{name: "Germany", country: "DE", year: 2025},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService()

			resp, err := svc.GetHolidays(tc.country, tc.year)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Country != tc.country {
				t.Errorf("expected country %q, got %q", tc.country, resp.Country)
			}
			if resp.Year != tc.year {
				t.Errorf("expected year %d, got %d", tc.year, resp.Year)
			}
			if len(resp.Holidays) == 0 {
				t.Error("expected non-empty holidays list")
			}
		})
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
			continue
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
	if !strings.Contains(err.Error(), "no holidays found") {
		t.Errorf("expected error message to contain 'no holidays found', got %q", err.Error())
	}
}
