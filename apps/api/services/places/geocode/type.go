package geocode

// GeocodeResponse is returned for address-to-coordinates lookups.
type GeocodeResponse struct { //nolint:revive // established public API type name
	Address string  `json:"address"`
	City    string  `json:"city"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func (GeocodeResponse) IsData() {}

// ReverseGeocodeResponse is returned for coordinates-to-address lookups.
type ReverseGeocodeResponse struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
	City    string  `json:"city"`
	Country string  `json:"country"`
}

func (ReverseGeocodeResponse) IsData() {}
