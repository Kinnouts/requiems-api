package qr

// Request holds the query parameters for the QR code endpoint.
type Request struct {
	Data   string `query:"data"   validate:"required"`
	Size   int    `query:"size"   validate:"min=50,max=1000"`
	Format string `query:"format" validate:"omitempty,oneof=png base64"`
}

// Base64Response is the JSON response payload when format=base64.
type Base64Response struct {
	Image  string `json:"image"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (Base64Response) IsData() {}
