package holidays

import (
	"errors"

	h "github.com/bobadilla-tech/holidays-per-country"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetHolidays(country string, year int) (HolidayList, error) {
	holidays := h.GetHolidays(country, year)
	if len(holidays) == 0 {
		return HolidayList{}, errors.New("no holidays found for the specified country and year")
	}

	holidayList := make([]Holiday, len(holidays))

	for i, holiday := range holidays {
		holidayList[i] = Holiday{
			Name: holiday.Name,
			Date: holiday.Date.Format("2006-01-02"),
		}
	}

	return HolidayList{
		Country:  country,
		Year:     year,
		Holidays: holidayList,
		Total:    len(holidayList),
	}, nil
}
