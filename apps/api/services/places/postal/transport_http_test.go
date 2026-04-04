package postal

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
		index: map[string]PostalCode{
			"US:10001": {
				PostalCode: "10001",
				City:       "New York City",
				State:      "New York",
				Country:    "US",
				Lat:        40.7484,
				Lon:        -73.9967,
			},
			"GB:SW1A1AA": {
				PostalCode: "SW1A1AA",
				City:       "London",
				State:      "England",
				Country:    "GB",
				Lat:        51.5014,
				Lon:        -0.1419,
			},
		},
	}
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestLookup_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/postal/10001?country=US", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[PostalCode]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.City != "New York City" {
		t.Errorf("expected city 'New York City', got %q", resp.Data.City)
	}
	if resp.Data.Lat == 0 {
		t.Error("expected non-zero latitude")
	}
}

func TestLookup_DefaultsToUS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/postal/10001", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 with default country=US, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLookup_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/postal/99999?country=US", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestLookup_NonUSCountry(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/postal/SW1A1AA?country=GB", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[PostalCode]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.Country != "GB" {
		t.Errorf("expected country 'GB', got %q", resp.Data.Country)
	}
}
