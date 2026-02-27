package timezone

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/timezone", func(w http.ResponseWriter, r *http.Request) {
		req := Request{}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		q := r.URL.Query()
		hasCoords := q.Has("lat") && q.Has("lon")
		hasCity := q.Has("city")

		if !hasCoords && !hasCity {
			httpx.Error(w, http.StatusBadRequest, "bad_request",
				"provide either 'city' or both 'lat' and 'lon' query parameters")
			return
		}

		var (
			info *Info
			err  error
		)

		if hasCity {
			info, err = svc.GetTimezoneByCity(req.City)
		} else {
			info, err = svc.GetTimezoneByCoords(req.Lat, req.Lon)
		}

		if err != nil {
			httpx.Error(w, http.StatusNotFound, "not_found", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, *info)
	})
}
