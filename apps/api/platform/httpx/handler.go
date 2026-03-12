package httpx

import (
	"context"
	"errors"
	"net/http"
)

// Handle wraps an endpoint function with automatic JSON binding, validation,
// and structured error responses. The request body is capped at 1 MB.
//
// Req is decoded from the JSON request body and validated via struct tags.
// Res must implement Data (add an IsData() method to your response type).
//
// Error mapping:
//   - *ValidationFailure  → 422 with {"error":"validation_failed","fields":[...]}
//   - *AppError           → AppError.Status with {"error":Code,"message":Message}
//   - any other error     → 500 with {"error":"internal_error"}
//
// Usage:
//
//	router.Post("/check", httpx.Handle(
//	    func(ctx context.Context, req CheckEmailRequest) (CheckEmailResponse, error) {
//	        return svc.CheckEmail(req.Email), nil
//	    },
//	))
func Handle[Req any, Res Data](
	fn func(ctx context.Context, req Req) (Res, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit

		var req Req

		if err := BindAndValidate(r, &req); err != nil {
			if vf, ok := errors.AsType[*ValidationFailure](err); ok {
				writeValidationError(w, vf.Fields)
				return
			}

			Error(w, http.StatusBadRequest, "bad_request", cleanDecodeError(err))
			return
		}

		res, err := fn(r.Context(), req)

		if err != nil {
			if ae, ok := errors.AsType[*AppError](err); ok {
				Error(w, ae.Status, ae.Code, ae.Message)
				return
			}

			Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}

		JSON(w, http.StatusOK, res)
	}
}

// cleanDecodeError returns a safe, human-readable message for JSON decode
// errors, hiding internal implementation details from the client.
func cleanDecodeError(err error) string {
	if err == nil {
		return ""
	}

	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return "request body too large (max 1MB)"
	}

	return "invalid request body"
}
