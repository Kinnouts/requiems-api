package asn

// IPAddressASNResponse is the JSON payload returned by the ASN lookup endpoint.
type IPAddressASNResponse struct {
	IP     string `json:"ip"`
	ASN    string `json:"asn"`
	Org    string `json:"org"`
	ISP    string `json:"isp"`
	Domain string `json:"domain"`
	Route  string `json:"route"`
	Type   string `json:"type"`
}

func (IPAddressASNResponse) IsData() {}
