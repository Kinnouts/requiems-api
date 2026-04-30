package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"requiems-api/platform/httpx"
)

type testData struct {
	Value string `json:"value"`
}

func (testData) IsData() {}

func TestJSON_WritesSuccessEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	httpx.JSON(w, http.StatusCreated, testData{Value: "hello"})

	if w.Code != http.StatusCreated {
		t.Errorf("status: want 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: want application/json, got %q", ct)
	}

	var resp httpx.Response[testData]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.Value != "hello" {
		t.Errorf("data.value: want %q, got %q", "hello", resp.Data.Value)
	}
	if resp.Metadata.Timestamp == "" {
		t.Error("metadata.timestamp must not be empty")
	}
}

func TestError_WritesErrorEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	httpx.Error(w, http.StatusNotFound, "not_found", "resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: want application/json, got %q", ct)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error != "not_found" {
		t.Errorf("error: want %q, got %q", "not_found", resp.Error)
	}
	if resp.Message != "resource not found" {
		t.Errorf("message: want %q, got %q", "resource not found", resp.Message)
	}
	if resp.Metadata.Timestamp == "" {
		t.Error("metadata.timestamp must not be empty")
	}
}
