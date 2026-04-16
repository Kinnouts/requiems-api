package qr

// Request holds the query parameters for the QR code endpoints.
type Request struct {
	Data     string `query:"data"     validate:"required"`
	Size     int    `query:"size"     validate:"min=50,max=1000"`
	Recovery string `query:"recovery" validate:"omitempty,oneof=low medium high highest"`
}

// Base64Response is the JSON response payload returned by GET /qr/base64.
type Base64Response struct {
	Image  string `json:"image"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (Base64Response) IsData() {}
