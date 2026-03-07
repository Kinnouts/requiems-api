package timezone

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter(t *testing.T) chi.Router {
	t.Helper()
	svc, err := NewService()
	if err != nil {
		t.Fatalf("failed to create timezone service: %v", err)
	}
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestTimezone_ByCoords(t *testing.T) {
	r := setupRouter(t)

	// London coordinates
	req := httptest.NewRequest(http.MethodGet, "/timezone?lat=51.5&lon=-0.1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Info]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Timezone != "Europe/London" {
		t.Errorf("expected timezone 'Europe/London', got %q", resp.Data.Timezone)
	}
	if resp.Data.CurrentTime == "" {
		t.Error("expected non-empty current_time")
	}
	if resp.Data.Offset == "" {
		t.Error("expected non-empty offset")
	}
}

func TestTimezone_ByCity(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/timezone?city=Tokyo", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Info]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Timezone != "Asia/Tokyo" {
		t.Errorf("expected timezone 'Asia/Tokyo', got %q", resp.Data.Timezone)
	}
}

func TestTimezone_CityNotFound(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/timezone?city=Atlantis", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestTimezone_MissingParams(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/timezone", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTimezone_MissingLon(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/timezone?lat=51.5", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTimezone_MissingLat(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/timezone?lon=-0.1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTimezone_InvalidLatRange(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/timezone?lat=200&lon=0", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTimezone_NewYork(t *testing.T) {
	r := setupRouter(t)

	// New York City coordinates
	req := httptest.NewRequest(http.MethodGet, "/timezone?lat=40.7&lon=-74.0", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Info]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Timezone != "America/New_York" {
		t.Errorf("expected timezone 'America/New_York', got %q", resp.Data.Timezone)
	}
}

func TestTimezone_CityLookup_CaseInsensitive(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	info, err := svc.GetTimezoneByCity("TOKYO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Timezone != "Asia/Tokyo" {
		t.Errorf("expected 'Asia/Tokyo', got %q", info.Timezone)
	}
}

func TestWorldTime_ValidTimezone(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/time/America/New_York", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Info]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Timezone != "America/New_York" {
		t.Errorf("expected timezone 'America/New_York', got %q", resp.Data.Timezone)
	}
	if resp.Data.CurrentTime == "" {
		t.Error("expected non-empty current_time")
	}
	if resp.Data.Offset == "" {
		t.Error("expected non-empty offset")
	}
}

func TestWorldTime_UTCTimezone(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/time/UTC", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Info]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Timezone != "UTC" {
		t.Errorf("expected timezone 'UTC', got %q", resp.Data.Timezone)
	}
	if resp.Data.Offset != "+00:00" {
		t.Errorf("expected offset '+00:00', got %q", resp.Data.Offset)
	}
}

func TestWorldTime_InvalidTimezone(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/time/Fake/Timezone", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestWorldTime_AsiaKolkata(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/time/Asia/Kolkata", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Info]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Timezone != "Asia/Kolkata" {
		t.Errorf("expected timezone 'Asia/Kolkata', got %q", resp.Data.Timezone)
	}
	if resp.Data.Offset != "+05:30" {
		t.Errorf("expected offset '+05:30', got %q", resp.Data.Offset)
	}
}

func TestTimezone_OffsetFormat(t *testing.T) {
	tests := []struct {
		offsetSecs int
		expected   string
	}{
		{0, "+00:00"},
		{3600, "+01:00"},
		{-18000, "-05:00"},
		{19800, "+05:30"},  // India
		{20700, "+05:45"},  // Nepal
		{-34200, "-09:30"}, // Marquesas Islands
	}

	for _, tc := range tests {
		got := formatOffset(tc.offsetSecs)
		if got != tc.expected {
			t.Errorf("formatOffset(%d) = %q, want %q", tc.offsetSecs, got, tc.expected)
		}
	}
}
