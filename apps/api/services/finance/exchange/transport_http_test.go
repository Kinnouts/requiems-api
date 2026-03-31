package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

type stubFetcher struct {
	fn func(ctx context.Context, from, to string) (float64, time.Time, error)
}

func (s *stubFetcher) GetRate(ctx context.Context, from, to string) (float64, time.Time, error) {
	return s.fn(ctx, from, to)
}

func fixedRate(rate float64) *stubFetcher {
	ts := time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)
	return &stubFetcher{fn: func(_ context.Context, _, _ string) (float64, time.Time, error) {
		return rate, ts, nil
	}}
}

func errFetcher(err error) *stubFetcher {
	return &stubFetcher{fn: func(_ context.Context, _, _ string) (float64, time.Time, error) {
		return 0, time.Time{}, err
	}}
}

func setupRouter(f Fetcher) chi.Router {
	r := chi.NewRouter()
	registerExchangeRoutes(r, f)
	return r
}

// — /exchange-rate tests —

func TestExchangeRate_HappyPath(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?from=USD&to=EUR", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ExchangeRateResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.From != "USD" {
		t.Errorf("from: want USD, got %s", resp.Data.From)
	}
	if resp.Data.To != "EUR" {
		t.Errorf("to: want EUR, got %s", resp.Data.To)
	}
	if resp.Data.Rate != 0.92 {
		t.Errorf("rate: want 0.92, got %v", resp.Data.Rate)
	}
	if resp.Data.Timestamp == "" {
		t.Error("timestamp must not be empty")
	}
	if resp.Metadata.Timestamp == "" {
		t.Error("metadata.timestamp must not be empty")
	}
}

func TestExchangeRate_LowercaseCodes_Normalized(t *testing.T) {
	var gotFrom, gotTo string
	f := &stubFetcher{fn: func(_ context.Context, from, to string) (float64, time.Time, error) {
		gotFrom, gotTo = from, to
		return 0.92, time.Now(), nil
	}}
	r := setupRouter(f)

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?from=usd&to=eur", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if gotFrom != "USD" || gotTo != "EUR" {
		t.Errorf("expected uppercase codes, got from=%s to=%s", gotFrom, gotTo)
	}
}

func TestExchangeRate_MissingFrom_Returns400(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?to=EUR", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExchangeRate_MissingTo_Returns400(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?from=USD", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExchangeRate_InvalidCurrencyCode_Returns400(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?from=US&to=EUR", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for 2-char code, got %d", w.Code)
	}
}

func TestExchangeRate_UnknownCurrency_Returns422(t *testing.T) {
	appErr := &httpx.AppError{Status: http.StatusUnprocessableEntity, Code: "invalid_currency", Message: "unknown currency code: XYZ"}
	r := setupRouter(errFetcher(appErr))

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?from=XYZ&to=EUR", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}

	var errResp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Error != "invalid_currency" {
		t.Errorf("expected error code invalid_currency, got %s", errResp.Error)
	}
}

func TestExchangeRate_UpstreamError_Returns503(t *testing.T) {
	r := setupRouter(errFetcher(errors.New("connection refused")))

	req := httptest.NewRequest(http.MethodGet, "/exchange-rate?from=USD&to=EUR", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

// — /convert tests —

func TestConvert_HappyPath(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR&amount=100", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ConvertResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.Amount != 100 {
		t.Errorf("amount: want 100, got %v", resp.Data.Amount)
	}
	if resp.Data.Converted != 92.00 {
		t.Errorf("converted: want 92.00, got %v", resp.Data.Converted)
	}
	if resp.Data.Rate != 0.92 {
		t.Errorf("rate: want 0.92, got %v", resp.Data.Rate)
	}
}

func TestConvert_MissingAmount_Returns400(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestConvert_ZeroAmount_Returns400(t *testing.T) {
	r := setupRouter(fixedRate(0.92))

	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR&amount=0", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for amount=0, got %d", w.Code)
	}
}

func TestConvert_UnknownCurrency_Returns422(t *testing.T) {
	appErr := &httpx.AppError{Status: http.StatusUnprocessableEntity, Code: "invalid_currency", Message: "unknown currency code: XYZ"}
	r := setupRouter(errFetcher(appErr))

	req := httptest.NewRequest(http.MethodGet, "/convert?from=XYZ&to=EUR&amount=50", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestConvert_UpstreamError_Returns503(t *testing.T) {
	r := setupRouter(errFetcher(errors.New("timeout")))

	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR&amount=100", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestConvert_ConversionRounding(t *testing.T) {
	r := setupRouter(fixedRate(0.9205))

	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR&amount=100", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp httpx.Response[ConvertResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.Converted != 92.05 {
		t.Errorf("converted: want 92.05, got %v", resp.Data.Converted)
	}
}
