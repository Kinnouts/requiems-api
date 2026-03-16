package timezone

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/time/*", func(w http.ResponseWriter, r *http.Request) {
		tzName := chi.URLParam(r, "*")
		if tzName == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "timezone is required")
			return
		}

		info, err := svc.GetCurrentTime(tzName)
		if err != nil {
			httpx.Error(w, http.StatusNotFound, "not_found", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, *info)
	})

	r.Get("/timezone", handleGetTimezone(svc))
}

func handleGetTimezone(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, hasCity, hasCoords, err := parseTimezoneQuery(r)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		var info *Info
		if hasCity {
			info, err = svc.GetTimezoneByCity(req.City)
		} else if hasCoords {
			info, err = svc.GetTimezoneByCoords(req.Lat, req.Lon)
		}

		if err != nil {
			httpx.Error(w, http.StatusNotFound, "not_found", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, *info)
	}
}

func parseTimezoneQuery(r *http.Request) (req Request, hasCity, hasCoords bool, err error) {
	if err = httpx.BindQuery(r, &req); err != nil {
		return
	}

	q := r.URL.Query()
	hasCoords = q.Has("lat") && q.Has("lon")
	hasCity = q.Has("city")

	if !hasCoords && !hasCity {
		err = errors.New("provide either 'city' or both 'lat' and 'lon' query parameters")
	}

	return
}
