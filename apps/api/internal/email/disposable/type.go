package disposable

// CheckEmailRequest represents a request to check if a single email is disposable
type CheckEmailRequest struct {
	Email string `json:"email"`
}

// CheckEmailResponse represents the response for a single email check
type CheckEmailResponse struct {
	Email        string `json:"email"`
	IsDisposable bool   `json:"is_disposable"`
	Domain       string `json:"domain,omitempty"`
}

// BatchCheckRequest represents a request to check multiple emails
type BatchCheckRequest struct {
	Emails []string `json:"emails"`
}

// BatchCheckResponse represents the response for a batch email check
type BatchCheckResponse struct {
	Results []CheckEmailResponse `json:"results"`
	Total   int                  `json:"total"`
}

// DomainCheckResponse represents the response for a domain check
type DomainCheckResponse struct {
	Domain       string `json:"domain"`
	IsDisposable bool   `json:"is_disposable"`
}

// DomainsListResponse represents the response for listing all domains
type DomainsListResponse struct {
	Domains []string `json:"domains"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	HasMore bool     `json:"has_more"`
}

// StatsResponse represents statistics about disposable domains
type StatsResponse struct {
	TotalDomains int `json:"total_domains"`
}
