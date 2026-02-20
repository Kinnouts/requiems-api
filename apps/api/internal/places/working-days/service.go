package workingdays

import (
	"time"

	businessdayscalculator "github.com/bobadilla-tech/business-days-calculator"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// GetWorkingDays calculates the number of working days between two dates
// It considers weekends and public holidays based on the provided country and subdivision
func (s *Service) GetWorkingDays(from, to time.Time, country string, subdivision string) int {
	if country == "" {
		return businessdayscalculator.CountBusinessDays(from, to)
	}

	opts := businessdayscalculator.HolidayOptions{
		CountryCode: country,
		Subdivision: subdivision,
	}
	return businessdayscalculator.CountBusinessDaysWithHolidays(from, to, opts)
}
