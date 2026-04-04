package cryptocoin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"requiems-api/platform/httpx"
)

func TestGetPrice_ValidSymbol(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer srv.Close()

	svc := newServiceWithClient(srv.Client(), srv.URL)
	p, err := svc.GetPrice(context.Background(), "BTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Symbol != "BTC" {
		t.Errorf("expected symbol BTC, got %s", p.Symbol)
	}
	if p.Name != "Bitcoin" {
		t.Errorf("expected name Bitcoin, got %s", p.Name)
	}
	if p.PriceUSD != 42000.50 {
		t.Errorf("expected price 42000.50, got %f", p.PriceUSD)
	}
	if p.Change24h != 2.5 {
		t.Errorf("expected change 2.5, got %f", p.Change24h)
	}
}

func TestGetPrice_UnknownSymbol(t *testing.T) {
	svc := newServiceWithClient(http.DefaultClient, "http://unused")
	_, err := svc.GetPrice(context.Background(), "FAKE")
	if err == nil {
		t.Fatal("expected error for unknown symbol")
	}

	ae, ok := err.(*httpx.AppError)
	if !ok {
		t.Fatalf("expected *httpx.AppError, got %T", err)
	}
	if ae.Code != "unknown_symbol" {
		t.Errorf("expected code unknown_symbol, got %s", ae.Code)
	}
	if ae.Status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", ae.Status)
	}
}

func TestGetPrice_UpstreamError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	svc := newServiceWithClient(srv.Client(), srv.URL)
	_, err := svc.GetPrice(context.Background(), "BTC")
	if err == nil {
		t.Fatal("expected error for upstream 500")
	}

	ae, ok := err.(*httpx.AppError)
	if !ok {
		t.Fatalf("expected *httpx.AppError, got %T", err)
	}
	if ae.Code != "upstream_error" {
		t.Errorf("expected code upstream_error, got %s", ae.Code)
	}
}

func TestGetPrice_NoRedis_CallsUpstream(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		body := coinGeckoResponse{
			"bitcoin": {USD: 50000},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	svc := newServiceWithClient(srv.Client(), srv.URL)

	for i := 0; i < 2; i++ {
		if _, err := svc.GetPrice(context.Background(), "BTC"); err != nil {
			t.Fatalf("call %d failed: %v", i+1, err)
		}
	}

	if callCount != 2 {
		t.Errorf("expected 2 upstream calls (no Redis), got %d", callCount)
	}
}
