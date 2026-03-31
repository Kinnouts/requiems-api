package tech

import (
	"log"

	"github.com/bobadilla-tech/go-ip-intelligence/v2/ipi"
	"github.com/go-chi/chi/v5"

	"requiems-api/platform/config"
	"requiems-api/services/tech/barcode"
	"requiems-api/services/tech/domain"
	"requiems-api/services/tech/ip/asn"
	"requiems-api/services/tech/ip/info"
	"requiems-api/services/tech/ip/vpn"
	"requiems-api/services/tech/mx"
	"requiems-api/services/tech/password"
	"requiems-api/services/tech/phone"
	"requiems-api/services/tech/qr"
	"requiems-api/services/tech/useragent"
	"requiems-api/services/tech/whois"
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

	domainSvc := domain.NewService()
	domain.RegisterRoutes(r, domainSvc)

	ipiClient, err := ipi.New(
		ipi.WithDatabasePath(cfg.VPNDatabasePath),
		ipi.WithASNDatabasePath(cfg.VPNASNDatabasePath),
		ipi.WithCityDatabasePath(cfg.IPCityDatabasePath),
	)
	if err != nil {
		log.Printf("tech: failed to initialize ip intelligence client; ip routes disabled: %v", err)
	}

	vpnSvc := vpn.NewService(ipiClient)
	vpn.RegisterRoutes(r, vpnSvc)

	asnSvc := asn.NewService(ipiClient)
	asn.RegisterRoutes(r, asnSvc)

	infoSvc := info.NewService(ipiClient)
	info.RegisterRoutes(r, infoSvc)

	whoisSvc := whois.NewService()
	whois.RegisterRoutes(r, whoisSvc)
	mxSvc := mx.NewService()
	mx.RegisterRoutes(r, mxSvc)
}
