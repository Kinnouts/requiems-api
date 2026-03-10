package holidays

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
	svc := NewService()
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestHolidays_ValidRequest(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays?country=US&year=2025", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Country != "US" {
		t.Errorf("expected country 'US', got %q", resp.Data.Country)
	}
	if resp.Data.Year != 2025 {
		t.Errorf("expected year 2025, got %d", resp.Data.Year)
	}
	if len(resp.Data.Holidays) == 0 {
		t.Error("expected non-empty holidays list")
	}
}

func TestHolidays_UK(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays?country=GB&year=2025", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Country != "GB" {
		t.Errorf("expected country 'GB', got %q", resp.Data.Country)
	}
}

func TestHolidays_MissingCountry(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays?year=2025", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHolidays_MissingYear(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays?country=US", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHolidays_InvalidCountry(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays?country=INVALID&year=2025", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHolidays_InvalidYear(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays?country=US&year=0", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHolidays_NoParams(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/holidays", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
