package app

import (
	"net/http"

	"requiems-api/internal/platform/httpx"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}