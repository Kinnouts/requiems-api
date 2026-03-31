package inflation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// stubGetter implements Getter for transport tests. It returns a fixed result
// or a fixed error on every call, keeping tests DB-free and fast.
type stubGetter struct {
	result InflationResponse
	err    error
}

func (s *stubGetter) GetInflation(_ context.Context, countryCode string) (InflationResponse, error) {
	if s.err != nil {
		return InflationResponse{}, s.err
	}
	r := s.result
	r.Country = countryCode
	return r, nil
}

// setupRouter wires up a stub getter into a chi router for handler testing.
func setupRouter(g Getter) chi.Router {
	r := chi.NewRouter()
	registerInflationRoutes(r, g)
	return r
}

// ---- helper ----

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) httpx.Response[InflationResponse] {
	t.Helper()
	var resp httpx.Response[InflationResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

// ---- tests ----

func TestInflation_KnownCountry_Returns200(t *testing.T) {
	svc := &stubGetter{result: InflationResponse{
		Rate:       3.2,
		Period:     "2024",
		Historical: []HistoricalRate{{Period: "2023", Rate: 4.1}},
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/inflation?country=US", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.Data.Country != "US" {
		t.Errorf("expected country US, got %q", resp.Data.Country)
	}
	if resp.Metadata.Timestamp == "" {
		t.Error("expected metadata.timestamp to be set")
	}
}

func TestInflation_ResponseEnvelope(t *testing.T) {
	svc := &stubGetter{result: InflationResponse{Rate: 2.5, Period: "2024"}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/inflation?country=DE", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if _, ok := raw["data"]; !ok {
		t.Error("response must have a 'data' key")
	}
	if _, ok := raw["metadata"]; !ok {
		t.Error("response must have a 'metadata' key")
	}
}

func TestInflation_UnknownCountry_Returns404(t *testing.T) {
	svc := &stubGetter{err: &httpx.AppError{
		Status:  http.StatusNotFound,
		Code:    "not_found",
		Message: "no inflation data found for country",
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/inflation?country=XK", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInflation_MissingCountryParam_Returns400(t *testing.T) {
	svc := &stubGetter{}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/inflation", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInflation_InvalidCountryCode_Returns400(t *testing.T) {
	svc := &stubGetter{}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/inflation?country=ZZZ", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInflation_LowercaseCountryAccepted(t *testing.T) {
	svc := &stubGetter{result: InflationResponse{Rate: 1.5, Period: "2024"}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/inflation?country=us", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for lowercase country, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.Data.Country != "US" {
		t.Errorf("expected country US (uppercased), got %q", resp.Data.Country)
	}
}

func TestInflation_HistoricalFieldPresent(t *testing.T) {
	historical := []HistoricalRate{
		{Period: "2023", Rate: 4.1},
		{Period: "2022", Rate: 8.0},
	}
	svc := &stubGetter{result: InflationResponse{
		Rate:       3.2,
		Period:     "2024",
		Historical: historical,
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/inflation?country=US", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := decodeResponse(t, w)
	if len(resp.Data.Historical) != 2 {
		t.Errorf("expected 2 historical entries, got %d", len(resp.Data.Historical))
	}
}
