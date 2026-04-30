package holidays

type Request struct {
	Country string `query:"country" validate:"required,iso3166_1_alpha2"`
	Year    int    `query:"year" validate:"required,min=1"`
}

type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

func (Holiday) IsData() {}

type HolidayList struct {
	Country  string    `json:"country"`
	Year     int       `json:"year"`
	Holidays []Holiday `json:"holidays"`
	Total    int       `json:"total"`
}

func (HolidayList) IsData() {}
