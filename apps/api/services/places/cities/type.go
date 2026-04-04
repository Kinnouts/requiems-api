package cities

// City is the response returned for a city lookup.
type City struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	Population int64   `json:"population"`
	Timezone   string  `json:"timezone"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

func (City) IsData() {}
