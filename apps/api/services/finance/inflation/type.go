package inflation

import "context"

// HistoricalRate is a single year's inflation rate.
type HistoricalRate struct {
	Period string  `json:"period"`
	Rate   float64 `json:"rate"`
}

// InflationResponse is the response payload for GET /v1/finance/inflation.
type InflationResponse struct {
	Country    string           `json:"country"`
	Rate       float64          `json:"rate"`
	Period     string           `json:"period"`
	Historical []HistoricalRate `json:"historical"`
}

func (InflationResponse) IsData() {}

// Request holds the validated query parameters for the inflation endpoint.
type Request struct {
	Country string `query:"country" validate:"required,iso3166_1_alpha2"`
}

// Getter is the interface used by the HTTP transport layer, allowing transport
// tests to inject a stub without requiring a database connection.
type Getter interface {
	GetInflation(ctx context.Context, countryCode string) (InflationResponse, error)
}
