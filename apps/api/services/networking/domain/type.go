package domain

// MXRecord holds an MX record entry.
type MXRecord struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
}

// DNSRecords holds the DNS records for a domain.
type DNSRecords struct {
	A     []string   `json:"a"`
	AAAA  []string   `json:"aaaa"`
	MX    []MXRecord `json:"mx"`
	NS    []string   `json:"ns"`
	TXT   []string   `json:"txt"`
	CNAME string     `json:"cname,omitempty"`
}

// InfoResponse is the response for a domain info request.
type InfoResponse struct {
	Domain    string     `json:"domain"`
	Available bool       `json:"available"`
	DNS       DNSRecords `json:"dns"`
}

func (InfoResponse) IsData() {}
