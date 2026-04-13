package httpx_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"requiems-api/platform/httpx"
)

// handleReq / handleRes are minimal request and response types used across
// Handle and HandleBatch tests.
type handleReq struct {
	Name string `json:"name" validate:"required"`
}

type handleRes struct {
	Greeting string `json:"greeting"`
}

func (handleRes) IsData() {}

func TestHandle_HappyPath(t *testing.T) {
	h := httpx.Handle(func(_ context.Context, req handleReq) (handleRes, error) {
		return handleRes{Greeting: "hello " + req.Name}, nil
	})

	body := strings.NewReader(`{"name":"world"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", w.Code)
	}

	var resp httpx.Response[handleRes]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.Greeting != "hello world" {
		t.Errorf("greeting: want %q, got %q", "hello world", resp.Data.Greeting)
	}
}

func TestHandle_MalformedJSON_Returns400(t *testing.T) {
	h := httpx.Handle(func(_ context.Context, req handleReq) (handleRes, error) {
		return handleRes{}, nil
	})

	body := strings.NewReader(`not-json`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

func TestHandle_ValidationFailure_Returns422(t *testing.T) {
	h := httpx.Handle(func(_ context.Context, req handleReq) (handleRes, error) {
		return handleRes{}, nil
	})

	// name is required but absent
	body := strings.NewReader(`{}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: want 422, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Error != "validation_failed" {
		t.Errorf("error: want validation_failed, got %q", resp.Error)
	}
	if len(resp.Fields) == 0 {
		t.Error("expected at least one field error")
	}
}

func TestHandle_AppError_MapsStatus(t *testing.T) {
	h := httpx.Handle(func(_ context.Context, req handleReq) (handleRes, error) {
		return handleRes{}, &httpx.AppError{Status: http.StatusTeapot, Code: "im_a_teapot", Message: "brew something"}
	})

	body := strings.NewReader(`{"name":"x"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusTeapot {
		t.Errorf("status: want 418, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Error != "im_a_teapot" {
		t.Errorf("error: want im_a_teapot, got %q", resp.Error)
	}
}

func TestHandle_InternalError_Returns500(t *testing.T) {
	h := httpx.Handle(func(_ context.Context, req handleReq) (handleRes, error) {
		return handleRes{}, errors.New("boom")
	})

	body := strings.NewReader(`{"name":"x"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Error != "internal_error" {
		t.Errorf("error: want internal_error, got %q", resp.Error)
	}
}

func TestHandleBatch_HappyPath_SetsUsageCountHeader(t *testing.T) {
	h := httpx.HandleBatch(func(_ context.Context, req handleReq) (handleRes, int, error) {
		return handleRes{Greeting: "hi " + req.Name}, 5, nil
	})

	body := strings.NewReader(`{"name":"batch"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", w.Code)
	}
	if got := w.Header().Get("X-Usage-Count"); got != "5" {
		t.Errorf("X-Usage-Count: want 5, got %q", got)
	}
}

func TestHandleBatch_ValidationFailure_Returns422(t *testing.T) {
	h := httpx.HandleBatch(func(_ context.Context, req handleReq) (handleRes, int, error) {
		return handleRes{}, 0, nil
	})

	body := strings.NewReader(`{}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: want 422, got %d", w.Code)
	}
}

func TestHandleBatch_AppError_MapsStatus(t *testing.T) {
	h := httpx.HandleBatch(func(_ context.Context, req handleReq) (handleRes, int, error) {
		return handleRes{}, 0, &httpx.AppError{Status: http.StatusBadGateway, Code: "bad_gateway", Message: "upstream down"}
	})

	body := strings.NewReader(`{"name":"x"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status: want 502, got %d", w.Code)
	}
}

func TestHandleBatch_InternalError_Returns500(t *testing.T) {
	h := httpx.HandleBatch(func(_ context.Context, req handleReq) (handleRes, int, error) {
		return handleRes{}, 0, errors.New("unexpected")
	})

	body := strings.NewReader(`{"name":"x"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", w.Code)
	}
}
