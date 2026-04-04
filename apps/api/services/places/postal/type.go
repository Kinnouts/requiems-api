package postal

// PostalCode is the response returned for a postal code lookup.
type PostalCode struct {
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

func (PostalCode) IsData() {}
