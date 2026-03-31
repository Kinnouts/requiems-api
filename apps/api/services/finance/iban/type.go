package iban

import "context"

// ParseResponse is the response payload for GET /v1/finance/iban/{iban}.
type ParseResponse struct {
	IBAN     string `json:"iban"`
	Valid    bool   `json:"valid"`
	Country  string `json:"country"`
	BankCode string `json:"bank_code"`
	Account  string `json:"account"`
}

func (ParseResponse) IsData() {}

// Validator is the interface used by the HTTP transport layer, allowing
// transport tests to inject a stub without requiring a database connection.
type Validator interface {
	Parse(ctx context.Context, raw string) (ParseResponse, error)
}
