package holidays

type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

type Request struct {
	Country string `query:"country" validate:"required,iso3166_1_alpha2"`
	Year    int    `query:"year" validate:"required,min=1"`
}

type Response struct {
	Country  string    `json:"country"`
	Year     int       `json:"year"`
	Holidays []Holiday `json:"holidays"`
}

func (Response) IsData() {}
