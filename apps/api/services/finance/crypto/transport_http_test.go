package cryptocoin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter(svc *Service) chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestCrypto_GetPrice_ValidSymbol(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := coinGeckoResponse{
			"bitcoin": {
				USD:          42000.50,
				USD24hChange: 2.5,
				USDMarketCap: 820000000000,
				USD24hVol:    25000000000,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}))
	defer upstream.Close()

	svc := newServiceWithClient(upstream.Client(), upstream.URL)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/crypto/BTC", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Price]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.Symbol != "BTC" {
		t.Errorf("expected symbol BTC, got %s", resp.Data.Symbol)
	}
	if resp.Data.Name != "Bitcoin" {
		t.Errorf("expected name Bitcoin, got %s", resp.Data.Name)
	}
	if resp.Data.PriceUSD != 42000.50 {
		t.Errorf("expected price 42000.50, got %f", resp.Data.PriceUSD)
	}
}

func TestCrypto_GetPrice_UppercaseNormalization(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := coinGeckoResponse{
			"bitcoin": {USD: 42000.50},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}))
	defer upstream.Close()

	svc := newServiceWithClient(upstream.Client(), upstream.URL)
	r := setupRouter(svc)

	// lowercase symbol should be normalized
	req := httptest.NewRequest(http.MethodGet, "/crypto/btc", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCrypto_GetPrice_UnknownSymbol(t *testing.T) {
	svc := newServiceWithClient(http.DefaultClient, "http://unused")
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/crypto/FAKE", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCrypto_GetPrice_UpstreamDown(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer upstream.Close()

	svc := newServiceWithClient(upstream.Client(), upstream.URL)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/crypto/ETH", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}
