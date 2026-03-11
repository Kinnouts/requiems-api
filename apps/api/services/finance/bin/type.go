package bin

import "context"

// LookupResponse is the response payload for GET /v1/finance/bin/:bin.
type LookupResponse struct {
	BIN         string  `json:"bin"`
	Scheme      string  `json:"scheme"`
	CardType    string  `json:"card_type"`
	CardLevel   string  `json:"card_level"`
	IssuerName  string  `json:"issuer_name"`
	IssuerURL   string  `json:"issuer_url"`
	IssuerPhone string  `json:"issuer_phone"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	Prepaid     bool    `json:"prepaid"`
	Luhn        bool    `json:"luhn"`
	Confidence  float64 `json:"confidence"`
}

func (LookupResponse) IsData() {}

// Looker is the interface used by the HTTP transport layer. Using an interface
// rather than a concrete *Service allows transport tests to inject a stub
// without requiring a database.
type Looker interface {
	Lookup(ctx context.Context, bin string) (LookupResponse, error)
}
