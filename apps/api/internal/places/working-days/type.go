package workingdays

import "time"

// WorkingDaysRequest holds the query parameters for the working days endpoint.
// Defaults should be set before calling httpx.BindQuery.
type WorkingDaysRequest struct {
	From    time.Time `query:"from" validate:"required"`
	To      time.Time `query:"to" validate:"required,gtfield=From"`
	Country string    `query:"country" validate:"omitempty,iso3166_1_alpha2"`
}

// WorkingDays represents the response for working days calculation
type WorkingDays struct {
	WorkingDays int    `json:"workingDays"`
	From        string `json:"from"`
	To          string `json:"to"`
	Country     string `json:"country,omitempty"`
}

func (WorkingDays) IsData() {}
