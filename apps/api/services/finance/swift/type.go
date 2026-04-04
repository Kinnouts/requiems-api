package swift

import "context"

// LookupResponse is the response payload for GET /v1/finance/swift/:code.
type LookupResponse struct {
	SwiftCode    string `json:"swift_code"`
	BankCode     string `json:"bank_code"`
	CountryCode  string `json:"country_code"`
	LocationCode string `json:"location_code"`
	BranchCode   string `json:"branch_code"`
	BankName     string `json:"bank_name"`
	City         string `json:"city"`
	CountryName  string `json:"country_name"`
	IsPrimary    bool   `json:"is_primary"`
}

func (LookupResponse) IsData() {}

// Looker is the interface used by the HTTP transport layer. Using an interface
// rather than a concrete *Service allows transport tests to inject a stub
// without requiring a database.
type Looker interface {
	Lookup(ctx context.Context, code string) (LookupResponse, error)
}
