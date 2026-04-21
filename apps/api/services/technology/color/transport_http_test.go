package color

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}

func TestColor_HappyPath_HexToRGB(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/color?from=hex&to=rgb&value=%23ffffff",
		nil,
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("expected application/json, got %s", ct)
	}

	var res httpx.Response[Response]
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Data.Input != "#ffffff" {
		t.Errorf("expected input #ffffff, got %s", res.Data.Input)
	}

	if !strings.Contains(res.Data.Result, "rgb") {
		t.Errorf("expected RGB result, got %s", res.Data.Result)
	}
}

func TestColor_MissingParam(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/color?from=hex&to=rgb",
		nil,
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String()) //Si falta un parámetro obligatorio, la API debe rechazar la request
	}
}

func TestColor_InvalidValue(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/color?from=hex&to=rgb&value=invalid",
		nil,
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestColor_ServiceError(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/color?from=hex&to=rgb&value=notahex",
		nil,
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected error status, got %d: %s", w.Code, w.Body.String())
	}

	ct := w.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("expected application/json, got %s", ct)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error != "invalid_color" {
		t.Fatalf("expected error code invalid_color, got %q", resp.Error)
	}
}
