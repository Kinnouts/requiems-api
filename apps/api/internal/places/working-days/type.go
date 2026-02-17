package workingdays

// WorkingDays represents the response for working days calculation
type WorkingDays struct {
	WorkingDays int    `json:"workingDays"`
	From        string `json:"from"`
	To          string `json:"to"`
	Country     string `json:"country,omitempty"`
}
