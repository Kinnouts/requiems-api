package barcode

// Request holds the query parameters for the barcode endpoints.
type Request struct {
	Data string `query:"data" validate:"required"`
	Type string `query:"type" validate:"required,oneof=code128 code93 code39 ean8 ean13"`
}

// Base64Response is the JSON response payload returned by GET /barcode/base64.
type Base64Response struct {
	Image  string `json:"image"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (Base64Response) IsData() {}
