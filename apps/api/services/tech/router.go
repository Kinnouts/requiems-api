package tech

import (
	"log"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/config"
	"requiems-api/services/tech/barcode"
	"requiems-api/services/tech/ip/vpn"
	"requiems-api/services/tech/password"
	"requiems-api/services/tech/phone"
	"requiems-api/services/tech/qr"
	"requiems-api/services/tech/useragent"
)

func RegisterRoutes(r chi.Router, cfg config.Config) {
	phoneSvc := phone.NewService()
	phone.RegisterRoutes(r, phoneSvc)

	passwordSvc := password.NewService()
	password.RegisterRoutes(r, passwordSvc)

	uaSvc := useragent.NewService()
	useragent.RegisterRoutes(r, uaSvc)

	qrSvc := qr.NewService()
	qr.RegisterRoutes(r, qrSvc)

	barcodeSvc := barcode.NewService()
	barcode.RegisterRoutes(r, barcodeSvc)

	vpnSvc, err := vpn.NewService(cfg.VPNDatabasePath, cfg.VPNASNDatabasePath)
	if err != nil {
		log.Fatalf("tech: failed to initialize vpn service: %v", err)
	}
	vpn.RegisterRoutes(r, vpnSvc)
}
