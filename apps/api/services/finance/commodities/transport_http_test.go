package commodities

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// stubGetterHTTP implements Getter for transport tests.
type stubGetterHTTP struct {
	result CommodityPrice
	err    error
}

func (s *stubGetterHTTP) Get(_ context.Context, slug string) (CommodityPrice, error) {
	if s.err != nil {
		return CommodityPrice{}, s.err
	}
	r := s.result
	r.Commodity = slug
	return r, nil
}

func setupRouter(g Getter) chi.Router {
	r := chi.NewRouter()
	registerCommodityRoutes(r, g)
	return r
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) httpx.Response[CommodityPrice] {
	t.Helper()
	var resp httpx.Response[CommodityPrice]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

// ---- tests ----

func TestCommodity_KnownSlug_Returns200(t *testing.T) {
	svc := &stubGetterHTTP{result: CommodityPrice{
		Name:      "Gold",
		Price:     2386.33,
		Unit:      "oz",
		Currency:  "USD",
		Change24h: 23.01,
		Historical: []HistoricalPrice{
			{Period: "2023", Price: 1940.54},
		},
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/commodities/gold", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.Data.Commodity != "gold" {
		t.Errorf("expected commodity gold, got %q", resp.Data.Commodity)
	}
	if resp.Data.Price != 2386.33 {
		t.Errorf("expected price 2386.33, got %v", resp.Data.Price)
	}
	if resp.Metadata.Timestamp == "" {
		t.Error("expected metadata.timestamp to be set")
	}
}

func TestCommodity_ResponseEnvelope(t *testing.T) {
	svc := &stubGetterHTTP{result: CommodityPrice{Price: 2386.33}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/commodities/gold", http.NoBody)
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

func TestCommodity_UnknownSlug_Returns404(t *testing.T) {
	svc := &stubGetterHTTP{err: &httpx.AppError{
		Status:  http.StatusNotFound,
		Code:    "not_found",
		Message: "commodity not found",
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/commodities/unobtainium", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCommodity_InternalError_Returns500(t *testing.T) {
	svc := &stubGetterHTTP{err: &httpx.AppError{
		Status:  http.StatusInternalServerError,
		Code:    "internal_error",
		Message: "internal server error",
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/commodities/gold", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCommodity_HistoricalFieldPresent(t *testing.T) {
	svc := &stubGetterHTTP{result: CommodityPrice{
		Price: 2386.33,
		Historical: []HistoricalPrice{
			{Period: "2023", Price: 1940.54},
			{Period: "2022", Price: 1800.12},
		},
	}}

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/commodities/gold", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := decodeResponse(t, w)
	if len(resp.Data.Historical) != 2 {
		t.Errorf("expected 2 historical entries, got %d", len(resp.Data.Historical))
	}
}
