package mx

// Record represents a single MX record entry.
type Record struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
}

// LookupResponse is the JSON payload returned by the MX lookup endpoint.
type LookupResponse struct {
	Domain  string     `json:"domain"`
	Records []Record `json:"records"`
}

func (LookupResponse) IsData() {}
