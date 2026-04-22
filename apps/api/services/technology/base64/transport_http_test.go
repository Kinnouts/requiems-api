package base64

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}
func TestBase64_Decode_OK(t *testing.T) {
	r := setupRouter()

	body := `{"value":"SGVsbG8="}`
	req := httptest.NewRequest(http.MethodPost, "/base64/decode", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var res httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}
}
func TestBase64_Encode_OK(t *testing.T) {
	r := setupRouter()

	body := `{"value" : "Hello"}`
	req := httptest.NewRequest(http.MethodPost, "/base64/encode" , bytes.NewBufferString(body))

	w := httptest.NewRecorder()

	r.ServeHTTP(w,req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if ct :=w.Header().Get("Content-Type")  ; ct !=  "application/json" {
			t.Errorf("Expected application/json , got  %s",ct)
	}

	var res httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if res.Data.Result == "" {
		t.Error("expected non-empty result")
	}
}

func TestBase64_Encode_MissingValue(t *testing.T) {
	r := setupRouter()

	body := `{"variant":"standard"}`
	req := httptest.NewRequest(http.MethodPost, "/base64/encode", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
	t.Errorf("expected application/json, got %s", ct)
	}

	var res httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if res.Error != "validation_error" && res.Error == "" {
		t.Error("expected error.code")
	}

	
}

func TestBase64_Encode_InvalidVariant(t *testing.T) {
	r := setupRouter()

	body := `{"value":"hello","variant":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/base64/encode", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
	t.Errorf("expected application/json, got %s", ct)
	}
	
	var res httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if res.Error == "" {
		t.Error("expected error.code")
	}

}

func TestBase64_Decode_InvalidBase64(t *testing.T) {
	r := setupRouter()

	body := `{"value":"%%%"}`
	req := httptest.NewRequest(http.MethodPost, "/base64/decode", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var res httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if res.Error == "" {
		t.Error("expected error.code")
	}

}
