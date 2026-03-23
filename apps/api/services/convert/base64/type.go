package base64

// EncodeRequest is the request body for the Base64 encode endpoint.
type EncodeRequest struct {
	Value   string `json:"value"   validate:"required"`
	Variant string `json:"variant" validate:"omitempty,oneof=standard url"`
}

// DecodeRequest is the request body for the Base64 decode endpoint.
type DecodeRequest struct {
	Value   string `json:"value"   validate:"required"`
	Variant string `json:"variant" validate:"omitempty,oneof=standard url"`
}

// Result is the response returned by both the encode and decode endpoints.
type Result struct {
	Result string `json:"result"`
}

func (Result) IsData() {}
