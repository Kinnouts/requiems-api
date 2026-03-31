package exchange

// ExchangeRateRequest holds the validated query parameters for GET /exchange-rate.
type ExchangeRateRequest struct {
	From string `query:"from" validate:"required,len=3,alpha"`
	To   string `query:"to"   validate:"required,len=3,alpha"`
}

// ConvertRequest holds the validated query parameters for GET /convert.
type ConvertRequest struct {
	From   string  `query:"from"   validate:"required,len=3,alpha"`
	To     string  `query:"to"     validate:"required,len=3,alpha"`
	Amount float64 `query:"amount" validate:"required,gt=0"`
}

// ExchangeRateResponse is the response payload for GET /v1/finance/exchange-rate.
type ExchangeRateResponse struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Rate      float64 `json:"rate"`
	Timestamp string  `json:"timestamp"`
}

func (ExchangeRateResponse) IsData() {}

// ConvertResponse is the response payload for GET /v1/finance/convert.
type ConvertResponse struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Rate      float64 `json:"rate"`
	Amount    float64 `json:"amount"`
	Converted float64 `json:"converted"`
	Timestamp string  `json:"timestamp"`
}

func (ConvertResponse) IsData() {}
