package info

// LookupResponse is the JSON payload returned by the IP geolocation endpoint.
type LookupResponse struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	IsVPN       bool   `json:"is_vpn"`
}

func (LookupResponse) IsData() {}
