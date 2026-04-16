package email

// Request holds the JSON body for email validation.
type Request struct {
	Email string `json:"email" validate:"required"`
}

// EmailValidation is the full validation result for an email address.
type EmailValidation struct {
	Email       *string `json:"email"`
	Valid       bool    `json:"valid"`
	SyntaxValid bool    `json:"syntax_valid"`
	MxValid     bool    `json:"mx_valid"`
	Disposable  bool    `json:"disposable"`
	Normalized  *string `json:"normalized"`
	Domain      *string `json:"domain"`
	Suggestion  *string `json:"suggestion"`
}

func (EmailValidation) IsData() {}
