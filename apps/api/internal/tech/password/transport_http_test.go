package password

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}

func TestPassword_DefaultLength(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/password", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Password]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Length != 16 {
		t.Errorf("expected length 16, got %d", resp.Data.Length)
	}

	if len(resp.Data.Password) != 16 {
		t.Errorf("expected password of length 16, got %d", len(resp.Data.Password))
	}
}

func TestPassword_CustomLength(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/password?length=32", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Password]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Length != 32 {
		t.Errorf("expected length 32, got %d", resp.Data.Length)
	}

	if len(resp.Data.Password) != 32 {
		t.Errorf("expected password of length 32, got %d", len(resp.Data.Password))
	}
}

func TestPassword_AllCharsets(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/password?length=64&uppercase=true&numbers=true&symbols=true", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Password]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	pwd := resp.Data.Password

	hasLower := strings.ContainsAny(pwd, charsetLower)
	hasUpper := strings.ContainsAny(pwd, charsetUpper)
	hasDigit := strings.ContainsAny(pwd, charsetNumbers)
	hasSymbol := strings.ContainsAny(pwd, charsetSymbols)

	if !hasLower {
		t.Error("expected at least one lowercase letter")
	}

	if !hasUpper {
		t.Error("expected at least one uppercase letter")
	}

	if !hasDigit {
		t.Error("expected at least one digit")
	}

	if !hasSymbol {
		t.Error("expected at least one symbol")
	}

	if resp.Data.Strength != "strong" {
		t.Errorf("expected strength 'strong', got %q", resp.Data.Strength)
	}
}

func TestPassword_LengthTooShort(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/password?length=4", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestPassword_LengthTooLong(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/password?length=200", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestPassword_StrengthWeak(t *testing.T) {
	svc := NewService()

	result, err := svc.Generate(8, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Strength != "weak" {
		t.Errorf("expected strength 'weak', got %q", result.Strength)
	}
}

func TestPassword_StrengthMedium(t *testing.T) {
	svc := NewService()

	result, err := svc.Generate(8, true, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Strength != "medium" {
		t.Errorf("expected strength 'medium', got %q", result.Strength)
	}
}

func TestPassword_StrengthStrong(t *testing.T) {
	svc := NewService()

	result, err := svc.Generate(16, true, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Strength != "strong" {
		t.Errorf("expected strength 'strong', got %q", result.Strength)
	}
}

func TestPassword_OnlyLowercase(t *testing.T) {
	svc := NewService()

	result, err := svc.Generate(12, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, c := range result.Password {
		if !strings.ContainsRune(charsetLower, c) {
			t.Errorf("unexpected character %q in lowercase-only password", c)
		}
	}
}

func TestPassword_NoSymbolsWhenNotRequested(t *testing.T) {
	svc := NewService()

	result, err := svc.Generate(32, true, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.ContainsAny(result.Password, charsetSymbols) {
		t.Error("expected no symbols in password when symbols not requested")
	}
}
