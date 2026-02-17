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
// The country parameter is currently ignored (stub for future implementation)
func (s *Service) GetWorkingDays(from, to time.Time, country string) int {
	return businessdayscalculator.CalculateBusinessDays(from, to)
}
