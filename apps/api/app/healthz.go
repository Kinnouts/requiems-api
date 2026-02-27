package app

import (
	"net/http"

	"requiems-api/platform/httpx"
)

type healthzResponse struct {
	Status string `json:"status"`
}

func (healthzResponse) IsData() {}

func Healthz(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, healthzResponse{Status: "ok"})
}
