package color //nolint:revive // package name matches the service domain it implements

// Request holds the validated query parameters for GET /convert/color.
type Request struct {
	From  string `query:"from"  validate:"required,oneof=hex rgb hsl cmyk"`
	To    string `query:"to"    validate:"required,oneof=hex rgb hsl cmyk"`
	Value string `query:"value" validate:"required"`
}

// Formats holds the color expressed in every supported format.
type Formats struct {
	Hex  string `json:"hex"`
	RGB  string `json:"rgb"`
	HSL  string `json:"hsl"`
	CMYK string `json:"cmyk"`
}

// Response is the payload returned by GET /v1/convert/color.
type Response struct {
	Input   string  `json:"input"`
	Result  string  `json:"result"`
	Formats Formats `json:"formats"`
}

func (Response) IsData() {}
