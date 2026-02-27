package useragent

// ParseRequest holds the query parameters for the user agent parse endpoint.
type ParseRequest struct {
	UA string `query:"ua" validate:"required"`
}

// Result holds parsed user agent information.
type Result struct {
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	OS             string `json:"os"`
	OSVersion      string `json:"os_version"`
	Device         string `json:"device"`
	IsBot          bool   `json:"is_bot"`
}

func (Result) IsData() {}
