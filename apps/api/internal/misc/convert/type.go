package convert

// ConvertResponse is the response payload for a unit conversion request.
type ConvertResponse struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Input   float64 `json:"input"`
	Result  float64 `json:"result"`
	Formula string  `json:"formula"`
}
