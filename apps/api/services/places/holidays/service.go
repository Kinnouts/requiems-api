package holidays

import (
	h "github.com/bobadilla-tech/holidays-per-country"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetHolidays(req Request) (Response, error) {
	holidays := h.GetHolidays(req.Country, req.Year)

	holidayList := make([]Holiday, len(holidays))

	for i, holiday := range holidays {
		holidayList[i] = Holiday{
			Name: holiday.Name,
			Date: holiday.Date.Format("2006-01-02"),
		}
	}

	return Response{
		Country:  req.Country,
		Year:     req.Year,
		Holidays: holidayList,
	}, nil
}
