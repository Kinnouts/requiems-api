package useragent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}

func TestUserAgent_HappyPath(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/useragent?ua=Mozilla%2F5.0+%28Windows+NT+10.0%3B+Win64%3B+x64%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Chrome%2F120.0.0.0+Safari%2F537.36", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Browser != "Chrome" {
		t.Errorf("expected browser Chrome, got %q", resp.Data.Browser)
	}
	if resp.Data.Device != "desktop" {
		t.Errorf("expected device desktop, got %q", resp.Data.Device)
	}
	if resp.Data.IsBot {
		t.Error("expected is_bot false")
	}
}

func TestUserAgent_MissingUA(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/useragent", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Error != "bad_request" {
		t.Errorf("expected error code bad_request, got %q", resp.Error)
	}
}

func TestUserAgent_BotDetection(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/useragent?ua=Mozilla%2F5.0+%28compatible%3B+Googlebot%2F2.1%3B+%2Bhttp%3A%2F%2Fwww.google.com%2Fbot.html%29", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Data.IsBot {
		t.Error("expected is_bot true for Googlebot")
	}
	if resp.Data.Device != "bot" {
		t.Errorf("expected device bot, got %q", resp.Data.Device)
	}
}
