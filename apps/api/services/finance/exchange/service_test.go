package exchange

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"requiems-api/platform/httpx"
)

// fakeFrankfurter returns a handler that serves a Frankfurter-shaped response.
func fakeFrankfurter(base, target string, rate float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := frankfurterResponse{
			Base:  base,
			Date:  "2024-12-15",
			Rates: map[string]float64{target: rate},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func newTestService(handler http.Handler) (svc *Service, cleanup func()) {
	srv := httptest.NewServer(handler)
	svc = newServiceWithClient(nil, srv.Client(), srv.URL)
	return svc, srv.Close
}

func TestGetRate_CacheMiss_FetchesFromAPI(t *testing.T) {
	svc, cleanup := newTestService(fakeFrankfurter("USD", "EUR", 0.92))
	defer cleanup()

	rate, ts, err := svc.GetRate(t.Context(), "USD", "EUR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate != 0.92 {
		t.Errorf("rate: want 0.92, got %v", rate)
	}
	if ts.IsZero() {
		t.Error("timestamp must not be zero")
	}
}

func TestGetRate_DateParsedCorrectly(t *testing.T) {
	svc, cleanup := newTestService(fakeFrankfurter("USD", "GBP", 0.78))
	defer cleanup()

	_, ts, err := svc.GetRate(t.Context(), "USD", "GBP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Year() != 2024 || ts.Month() != 12 || ts.Day() != 15 {
		t.Errorf("expected date 2024-12-15, got %s", ts.Format("2006-01-02"))
	}
}

func TestGetRate_InvalidTargetCurrency_Returns422(t *testing.T) {
	// API returns empty rates map for an unknown target.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := frankfurterResponse{Base: "USD", Date: "2024-12-15", Rates: map[string]float64{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	svc, cleanup := newTestService(handler)
	defer cleanup()

	_, _, err := svc.GetRate(t.Context(), "USD", "XYZ")
	if err == nil {
		t.Fatal("expected error for unknown target currency")
	}

	var ae *httpx.AppError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *httpx.AppError, got %T", err)
	}
	if ae.Code != "invalid_currency" {
		t.Errorf("code: want invalid_currency, got %s", ae.Code)
	}
	if ae.Status != http.StatusUnprocessableEntity {
		t.Errorf("status: want 422, got %d", ae.Status)
	}
}

func TestGetRate_APIReturns404_Returns422(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	svc, cleanup := newTestService(handler)
	defer cleanup()

	_, _, err := svc.GetRate(t.Context(), "XYZ", "EUR")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}

	var ae *httpx.AppError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *httpx.AppError, got %T", err)
	}
	if ae.Code != "invalid_currency" {
		t.Errorf("code: want invalid_currency, got %s", ae.Code)
	}
}

func TestGetRate_APIReturns500_Returns503(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	svc, cleanup := newTestService(handler)
	defer cleanup()

	_, _, err := svc.GetRate(t.Context(), "USD", "EUR")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}

	var ae *httpx.AppError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *httpx.AppError, got %T", err)
	}
	if ae.Code != "upstream_error" {
		t.Errorf("code: want upstream_error, got %s", ae.Code)
	}
	if ae.Status != http.StatusServiceUnavailable {
		t.Errorf("status: want 503, got %d", ae.Status)
	}
}

func TestParseCache_RoundTrip(t *testing.T) {
	rate := 0.9205
	ts := time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)
	val := formatCacheValue(rate, ts)

	gotRate, gotTS, err := parseCache(val)
	if err != nil {
		t.Fatalf("parseCache: %v", err)
	}
	if gotRate != rate {
		t.Errorf("rate: want %v, got %v", rate, gotRate)
	}
	if !gotTS.Equal(ts) {
		t.Errorf("ts: want %v, got %v", ts, gotTS)
	}
}

func TestParseCache_InvalidFormat(t *testing.T) {
	_, _, err := parseCache("notvalid")
	if err == nil {
		t.Error("expected error for invalid cache value")
	}
}

func TestCacheKey_AlwaysUppercase(t *testing.T) {
	if cacheKey("usd", "eur") != "exchange:USD:EUR" {
		t.Error("cache key must be uppercase")
	}
	if cacheKey("USD", "EUR") != "exchange:USD:EUR" {
		t.Error("cache key must be uppercase")
	}
}
