package markdown

// Request is the input for the markdown-to-HTML conversion endpoint.
type Request struct {
	Markdown string `json:"markdown" validate:"required"`
	Sanitize bool   `json:"sanitize"`
}

// Response holds the converted HTML output.
type Response struct {
	HTML string `json:"html"`
}

func (Response) IsData() {}
