package crypto

// CryptoPrice is the response payload for GET /v1/finance/crypto/{symbol}.
type CryptoPrice struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	PriceUSD  float64 `json:"price_usd"`
	Change24h float64 `json:"change_24h"`
	MarketCap float64 `json:"market_cap"`
	Volume24h float64 `json:"volume_24h"`
}

func (CryptoPrice) IsData() {}
