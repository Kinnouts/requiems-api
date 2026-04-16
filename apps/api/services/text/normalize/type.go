package normalize

import (
	normalizer "github.com/bobadilla-tech/go-email-normalizer"
)

type EmailNormalizationRequest struct {
	Email string `json:"email" validate:"required"`
}

type EmailNormalization struct {
	Original   string              `json:"original"`
	Normalized string              `json:"normalized"`
	Local      string              `json:"local"`
	Domain     string              `json:"domain"`
	Changes    []normalizer.Change `json:"changes"`
}

func (EmailNormalization) IsData() {}
