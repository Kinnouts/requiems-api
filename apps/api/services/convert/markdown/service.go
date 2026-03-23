package markdown

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Convert renders markdown to HTML following the CommonMark spec.
// When sanitize is true, raw HTML blocks and inline HTML in the markdown input
// are stripped from the output instead of being passed through.
func (s *Service) Convert(markdown string, sanitize bool) (Response, error) {
	opts := []goldmark.Option{
		goldmark.WithExtensions(extension.GFM),
	}

	if !sanitize {
		opts = append(opts, goldmark.WithRendererOptions(html.WithUnsafe()))
	}

	md := goldmark.New(opts...)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return Response{}, err
	}

	return Response{HTML: strings.TrimRight(buf.String(), "\n")}, nil
}
