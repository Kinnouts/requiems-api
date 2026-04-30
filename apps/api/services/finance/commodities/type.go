package commodities

import "context"

// HistoricalPrice is a single year's average commodity price.
type HistoricalPrice struct {
	Period string  `json:"period"`
	Price  float64 `json:"price"`
}

// CommodityPrice is the response payload for GET /v1/finance/commodities/:commodity.
type CommodityPrice struct {
	Commodity  string            `json:"commodity"`
	Name       string            `json:"name"`
	Price      float64           `json:"price"`
	Unit       string            `json:"unit"`
	Currency   string            `json:"currency"`
	Change24h  float64           `json:"change_24h"`
	Historical []HistoricalPrice `json:"historical"`
}

func (CommodityPrice) IsData() {}

// Getter is the interface used by the HTTP transport layer, allowing transport
// tests to inject a stub without requiring a database connection.
type Getter interface {
	Get(ctx context.Context, slug string) (CommodityPrice, error)
}
