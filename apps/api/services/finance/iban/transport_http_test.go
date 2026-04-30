package iban

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// stubValidator implements Validator for transport tests. It returns a fixed
// result or a fixed error on every call, keeping tests DB-free and fast.
type stubValidator struct {
	result ParseResponse
	err    error
}

func (s *stubValidator) Parse(_ context.Context, raw string) (ParseResponse, error) {
	if s.err != nil {
		return ParseResponse{}, s.err
	}
	r := s.result
	r.IBAN = raw
	return r, nil
}

func setupRouter(v Validator) chi.Router {
	r := chi.NewRouter()
	registerIBANRoutes(r, v)
	return r
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) httpx.Response[ParseResponse] {
	t.Helper()
	var resp httpx.Response[ParseResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

// ---- tests ----

func TestIBAN_ValidDE_Returns200(t *testing.T) {
	svc := &stubValidator{result: ParseResponse{
		Valid:    true,
		Country:  "Germany",
		BankCode: "37040044",
		Account:  "0532013000",
	}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/DE89370400440532013000", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if !resp.Data.Valid {
		t.Error("expected valid = true")
	}
	if resp.Data.IBAN != "DE89370400440532013000" {
		t.Errorf("expected IBAN echoed, got %q", resp.Data.IBAN)
	}
	if resp.Data.Country != "Germany" {
		t.Errorf("expected country Germany, got %q", resp.Data.Country)
	}
	if resp.Data.BankCode != "37040044" {
		t.Errorf("expected bank_code 37040044, got %q", resp.Data.BankCode)
	}
	if resp.Data.Account != "0532013000" {
		t.Errorf("expected account 0532013000, got %q", resp.Data.Account)
	}
	if resp.Metadata.Timestamp == "" {
		t.Error("expected metadata.timestamp to be set")
	}
}

func TestIBAN_ResponseEnvelope(t *testing.T) {
	svc := &stubValidator{result: ParseResponse{Valid: true, Country: "Netherlands"}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/NL91ABNA0417164300", http.NoBody)
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

func TestIBAN_InvalidChecksum_Returns200WithValidFalse(t *testing.T) {
	svc := &stubValidator{result: ParseResponse{Valid: false, Country: "Germany"}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/DE00370400440532013000", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for invalid IBAN, got %d", w.Code)
	}

	resp := decodeResponse(t, w)
	if resp.Data.Valid {
		t.Error("expected valid = false for invalid IBAN")
	}
}

func TestIBAN_DBError_Returns500(t *testing.T) {
	svc := &stubValidator{err: errors.New("db unavailable")}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/DE89370400440532013000", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for DB error, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIBAN_UnknownCountry_Returns200(t *testing.T) {
	// IBAN from country not in DB — valid checksum, empty bank_code/account.
	svc := &stubValidator{result: ParseResponse{Valid: true}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/XX00TEST12345678", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestIBAN_AllResponseFieldsPresent(t *testing.T) {
	svc := &stubValidator{result: ParseResponse{
		Valid:    true,
		Country:  "Netherlands",
		BankCode: "ABNA",
		Account:  "0417164300",
	}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/NL91ABNA0417164300", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := decodeResponse(t, w)
	d := resp.Data

	checks := map[string]bool{
		"iban non-empty":      d.IBAN != "",
		"country non-empty":   d.Country != "",
		"bank_code non-empty": d.BankCode != "",
		"account non-empty":   d.Account != "",
		"valid is true":       d.Valid,
	}
	for name, ok := range checks {
		if !ok {
			t.Errorf("field check failed: %s", name)
		}
	}
}

func TestIBAN_GBParsing_Returns200(t *testing.T) {
	svc := &stubValidator{result: ParseResponse{
		Valid:    true,
		Country:  "United Kingdom",
		BankCode: "WEST",
		Account:  "98765432",
	}}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/iban/GB82WEST12345698765432", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.Data.BankCode != "WEST" {
		t.Errorf("expected bank_code WEST, got %q", resp.Data.BankCode)
	}
	if resp.Data.Account != "98765432" {
		t.Errorf("expected account 98765432, got %q", resp.Data.Account)
	}
}
