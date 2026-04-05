package numbase

// Result is the response returned by the base conversion endpoint.
type Result struct {
	Input  string `json:"input"`
	From   int    `json:"from"`
	To     int    `json:"to"`
	Result string `json:"result"`
}

func (Result) IsData() {}
