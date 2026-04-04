package cities

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	svc := &Service{
		index: map[string]City{
			"london": {
				Name:       "London",
				Country:    "GB",
				Population: 7556900,
				Timezone:   "Europe/London",
				Lat:        51.5085,
				Lon:        -0.1257,
			},
			"new york city": {
				Name:       "New York City",
				Country:    "US",
				Population: 8336817,
				Timezone:   "America/New_York",
				Lat:        40.7128,
				Lon:        -74.0060,
			},
		},
	}
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestFind_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cities/london", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[City]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.Name != "London" {
		t.Errorf("expected name 'London', got %q", resp.Data.Name)
	}
	if resp.Data.Country != "GB" {
		t.Errorf("expected country 'GB', got %q", resp.Data.Country)
	}
	if resp.Data.Timezone != "Europe/London" {
		t.Errorf("expected timezone 'Europe/London', got %q", resp.Data.Timezone)
	}
}

func TestFind_CaseInsensitive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cities/London", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for mixed-case lookup, got %d", w.Code)
	}
}

func TestFind_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cities/atlantis", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestFind_MultiWordCity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cities/new york city", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[City]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.Population == 0 {
		t.Error("expected non-zero population")
	}
}
