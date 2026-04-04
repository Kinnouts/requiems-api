package geocode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter(mockServer *httptest.Server) chi.Router {
	svc := NewService(mockServer.URL, mockServer.Client(), nil)
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestGeocode_HappyPath(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"lat":"38.8976763","lon":"-77.0365298","display_name":"White House, Washington, DC","address":{"city":"Washington","country_code":"us"}}]`)) //nolint:errcheck
	}))
	defer mock.Close()

	req := httptest.NewRequest(http.MethodGet, "/geocode?address=1600+Pennsylvania+Ave", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter(mock).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[GeocodeResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.City != "Washington" {
		t.Errorf("expected city 'Washington', got %q", resp.Data.City)
	}
	if resp.Data.Country != "US" {
		t.Errorf("expected country 'US', got %q", resp.Data.Country)
	}
	if resp.Data.Lat == 0 {
		t.Error("expected non-zero latitude")
	}
}

func TestGeocode_MissingAddress(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`[]`)) //nolint:errcheck
	}))
	defer mock.Close()

	req := httptest.NewRequest(http.MethodGet, "/geocode", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter(mock).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing address, got %d", w.Code)
	}
}

func TestGeocode_NoResults(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`)) //nolint:errcheck
	}))
	defer mock.Close()

	req := httptest.NewRequest(http.MethodGet, "/geocode?address=zzznoresultsxxx", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter(mock).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGeocode_UpstreamError(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mock.Close()

	req := httptest.NewRequest(http.MethodGet, "/geocode?address=anywhere", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter(mock).ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestReverseGeocode_HappyPath(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"display_name":"White House, Washington, DC","address":{"city":"Washington","country_code":"us"}}`)) //nolint:errcheck
	}))
	defer mock.Close()

	req := httptest.NewRequest(http.MethodGet, "/reverse-geocode?lat=38.8977&lon=-77.0365", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter(mock).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ReverseGeocodeResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.City != "Washington" {
		t.Errorf("expected city 'Washington', got %q", resp.Data.City)
	}
}

func TestReverseGeocode_MissingParams(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer mock.Close()

	req := httptest.NewRequest(http.MethodGet, "/reverse-geocode?lat=38.8977", http.NoBody)
	w := httptest.NewRecorder()
	setupRouter(mock).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing lon, got %d", w.Code)
	}
}
