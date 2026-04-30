package convformat

// Request is the input for the format conversion endpoint.
type Request struct {
	From    string `json:"from"    validate:"required,oneof=json yaml csv xml toml"`
	To      string `json:"to"      validate:"required,oneof=json yaml csv xml toml"`
	Content string `json:"content" validate:"required"`
}

// Response holds the converted output.
type Response struct {
	Result string `json:"result"`
}

func (Response) IsData() {}
