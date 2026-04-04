package swift

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"requiems-api/platform/httpx"
)

// Service provides SWIFT/BIC code lookup against the swift_codes PostgreSQL table.
type Service struct {
	db *pgxpool.Pool
}

// NewService creates a new Service backed by the given connection pool.
func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

// Lookup validates the raw SWIFT/BIC code and queries the database for bank
// metadata. 8-character codes are transparently expanded to 11 characters by
// appending "XXX" (primary office lookup).
func (s *Service) Lookup(ctx context.Context, raw string) (LookupResponse, error) {
	code, err := sanitizeSWIFT(raw)
	if err != nil {
		return LookupResponse{}, err
	}

	result, err := s.querySwift(ctx, code)
	if errors.Is(err, pgx.ErrNoRows) {
		return LookupResponse{}, &httpx.AppError{
			Status:  http.StatusNotFound,
			Code:    "not_found",
			Message: "SWIFT code not found",
		}
	}
	if err != nil {
		return LookupResponse{}, err
	}

	result.IsPrimary = result.BranchCode == "XXX"
	return result, nil
}

// querySwift executes a point-lookup against swift_codes by exact swift_code.
func (s *Service) querySwift(ctx context.Context, code string) (LookupResponse, error) {
	row := s.db.QueryRow(ctx, `
		SELECT
			swift_code, bank_code, country_code, location_code, branch_code,
			bank_name, city, country_name
		FROM swift_codes
		WHERE swift_code = $1
	`, code)

	var r LookupResponse
	err := row.Scan(
		&r.SwiftCode, &r.BankCode, &r.CountryCode, &r.LocationCode, &r.BranchCode,
		&r.BankName, &r.City, &r.CountryName,
	)
	return r, err
}

// sanitizeSWIFT normalises the raw input and validates it against the SWIFT/BIC
// format (ISO 9362). Returns the canonical 11-character code (appending "XXX"
// for 8-character primary office codes) or an *httpx.AppError for client errors.
func sanitizeSWIFT(raw string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(raw))

	if len(code) != 8 && len(code) != 11 {
		return "", &httpx.AppError{
			Status:  http.StatusBadRequest,
			Code:    "bad_request",
			Message: "SWIFT code must be 8 or 11 characters",
		}
	}

	// Positions 0–3: bank code — must be letters only.
	for i := range 4 {
		if code[i] < 'A' || code[i] > 'Z' {
			return "", &httpx.AppError{
				Status:  http.StatusBadRequest,
				Code:    "bad_request",
				Message: "bank code must be 4 letters (positions 1–4)",
			}
		}
	}

	// Positions 4–5: country code — must be letters only.
	for i := 4; i < 6; i++ {
		if code[i] < 'A' || code[i] > 'Z' {
			return "", &httpx.AppError{
				Status:  http.StatusBadRequest,
				Code:    "bad_request",
				Message: "country code must be 2 letters (positions 5–6)",
			}
		}
	}

	// Positions 6–7: location code — alphanumeric.
	for i := 6; i < 8; i++ {
		if !isAlphanumeric(code[i]) {
			return "", &httpx.AppError{
				Status:  http.StatusBadRequest,
				Code:    "bad_request",
				Message: "location code must be alphanumeric (positions 7–8)",
			}
		}
	}

	// Positions 8–10: branch code — alphanumeric (only present for 11-char codes).
	if len(code) == 11 {
		for i := 8; i < 11; i++ {
			if !isAlphanumeric(code[i]) {
				return "", &httpx.AppError{
					Status:  http.StatusBadRequest,
					Code:    "bad_request",
					Message: "branch code must be alphanumeric (positions 9–11)",
				}
			}
		}
	}

	if len(code) == 8 {
		code += "XXX"
	}

	return code, nil
}

// isAlphanumeric reports whether b is an ASCII letter or digit.
func isAlphanumeric(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
