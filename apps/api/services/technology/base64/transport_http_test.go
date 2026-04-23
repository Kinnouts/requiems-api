package base64 //nolint:revive //En Go, todos los archivos que tienen el mismo package se ven entre sí automáticamente. 

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"requiems-api/platform/httpx"
)

// setupRouter arma el router con el servicio real en memoria, sin levantar ningún servidor.
// Service no tiene dependencias externas (sin base de datos, sin HTTP)
func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}


// verificarJSON comprueba que la respuesta tenga Content-Type JSON y body JSON válido.
// Se llama en todos los casos porque el ticket exige verificar estas dos cosas.
func assertJSON(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type: got %q, want application/json", ct)
	}
	if !json.Valid(w.Body.Bytes()) {
		t.Errorf("body is not valid JSON: %s", w.Body.String())
	}
}

// ── /base64/encode ────────────────────────────────────────────────────────────
// TestEncode_HappyPath verifica que el endpoint responde correctamente
// cuando recibe un request completo y válido.
func TestEncode_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/encode",
		strings.NewReader(`{"value":"Hello, world!"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.Result != "SGVsbG8sIHdvcmxkIQ==" {
		t.Errorf("result: got %q, want %q", resp.Data.Result, "SGVsbG8sIHdvcmxkIQ==")
	}
}


// TestEncode_MissngValue verifica que el endpoint rechaza el request
// cuando no viene el campo obligatorio "value".
// El framework detecta esto automáticamente y devuelve 422
// antes de que el servicio llegue a ejecutarse.
func TestEncode_MissingValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/encode",
		strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)
}

// TestEncode_InvalidVAriant verifica que el endpoint rechaza el request
// cuando "variant" recibe un valor que no existe (solo acepta "standard" o "url").
// El framework lo detecta por el tag validate:"oneof=standard url" y devuelve 422.
func TestEncode_InvalidVariant(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/encode",
		strings.NewReader(`{"value":"Hello","variant":"invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)
}

// ── /base64/decode ────────────────────────────────────────────────────────────
// TestDecode_HAppyPath verifica que el endpoint decodifica correctamente
// cuando recibe un base64 válido.
func TestDecode_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/decode",
		strings.NewReader(`{"value":"SGVsbG8sIHdvcmxkIQ=="}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)

	var resp httpx.Response[Result] //empleando librería interna del proyecto para respuesta estandarizada
	if err := json.NewDecoder(w.Body).Decode(&resp); 
		err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.Result != "Hello, world!" {
		t.Errorf("result: got %q, want %q", resp.Data.Result, "Hello, world!")
	}
}

/// TestDecode_MissingValue verifica que el endpoint rechaza el request
// cuando no viene el campo obligatorio "value".
func TestDecode_MissingValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/decode",
		strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)
}

// TestDecode_InvalidVariant verifica que el endpoint rechaza el request
// cuando "variant" recibe un valor que no existe.
func TestDecode_InvalidVariant(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/decode",
		strings.NewReader(`{"value":"SGVsbG8=","variant":"invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)
}

// TestDecode_ServiceError verifica el caso donde el request pasa la validación
// (tiene "value", el "variant" vacío es válido), pero el servicio falla porque
// el string no es base64 real. A diferencia de los casos anteriores, acá el
// servicio SÍ se ejecuta y es él quien devuelve el error con código 422.
func TestDecode_ServiceError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/base64/decode",
		strings.NewReader(`{"value":"not-valid-base64!!!"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	assertJSON(t, w)
}
