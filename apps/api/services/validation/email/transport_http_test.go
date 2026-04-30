package email

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

func postValidate(body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/email", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, req)
	return w
}

func TestValidate_MissingEmail(t *testing.T) {
	w := postValidate(`{}`)

	// httpx.Handle returns 422 for failed struct validation.
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestValidate_InvalidSyntax(t *testing.T) {
	w := postValidate(`{"email":"notanemail"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Validation]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Valid {
		t.Error("expected Valid=false for invalid syntax")
	}
	if resp.Data.SyntaxValid {
		t.Error("expected SyntaxValid=false for invalid syntax")
	}
}

func TestValidate_ValidEmail(t *testing.T) {
	w := postValidate(`{"email":"user@gmail.com"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Validation]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Data.SyntaxValid {
		t.Error("expected SyntaxValid=true")
	}
	if *resp.Data.Email != "user@gmail.com" {
		t.Errorf("expected Email=user@gmail.com, got %q", *resp.Data.Email)
	}
	if *resp.Data.Domain != "gmail.com" {
		t.Errorf("expected Domain=gmail.com, got %q", *resp.Data.Domain)
	}
}

func TestValidate_SuggestionPresentInResponse(t *testing.T) {
	w := postValidate(`{"email":"user@gmial.com"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Decode into a raw map to verify "suggestion" is always present, even when null.
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	var data map[string]json.RawMessage
	if err := json.Unmarshal(raw["data"], &data); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if _, ok := data["suggestion"]; !ok {
		t.Error("expected 'suggestion' key to be present in response (even when null)")
	}

	var suggestion *string
	if err := json.Unmarshal(data["suggestion"], &suggestion); err != nil {
		t.Fatalf("failed to decode suggestion: %v", err)
	}
	if suggestion == nil {
		t.Fatal("expected non-nil suggestion for gmial.com")
	}
	if *suggestion != "gmail.com" {
		t.Errorf("expected suggestion=gmail.com, got %q", *suggestion)
	}
}

func TestValidate_SuggestionNullForKnownDomain(t *testing.T) {
	w := postValidate(`{"email":"user@gmail.com"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	var data map[string]json.RawMessage
	if err := json.Unmarshal(raw["data"], &data); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if _, ok := data["suggestion"]; !ok {
		t.Error("expected 'suggestion' key to always be present")
	}

	if string(data["suggestion"]) != "null" {
		t.Errorf("expected suggestion=null for gmail.com, got %s", data["suggestion"])
	}
}
