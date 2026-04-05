package numbase

// ConvertRequest holds the validated query parameters for the base conversion endpoint.
// Defaults should be set before calling httpx.BindQuery.
type ConvertRequest struct {
	From  int    `query:"from"  validate:"required,oneof=2 8 10 16"`
	To    int    `query:"to"    validate:"required,oneof=2 8 10 16"`
	Value string `query:"value" validate:"required"`
}

// Result is the response returned by the base conversion endpoint.
type Result struct {
	Input  string `json:"input"`
	From   int    `json:"from"`
	To     int    `json:"to"`
	Result string `json:"result"`
}

func (Result) IsData() {}
