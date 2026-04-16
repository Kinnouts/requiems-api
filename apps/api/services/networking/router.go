package networking

import (
	"log"

	"github.com/bobadilla-tech/go-ip-intelligence/v2/ipi"
	"github.com/go-chi/chi/v5"

	"requiems-api/platform/config"
	"requiems-api/services/networking/disposable"
	"requiems-api/services/networking/domain"
	"requiems-api/services/networking/ip/asn"
	"requiems-api/services/networking/ip/info"
	"requiems-api/services/networking/ip/vpn"
	"requiems-api/services/networking/mx"
	"requiems-api/services/networking/whois"
)

// RegisterRoutes wires all networking sub-services onto the given router.
func RegisterRoutes(r chi.Router, cfg config.Config) {
	disposableSvc := disposable.NewService()
	disposable.RegisterRoutes(r, disposableSvc)

	domainSvc := domain.NewService()
	domain.RegisterRoutes(r, domainSvc)

	whoisSvc := whois.NewService()
	whois.RegisterRoutes(r, whoisSvc)

	mxSvc := mx.NewService()
	mx.RegisterRoutes(r, mxSvc)

	ipiClient, err := ipi.New(
		ipi.WithDatabasePath(cfg.VPNDatabasePath),
		ipi.WithASNDatabasePath(cfg.VPNASNDatabasePath),
		ipi.WithCityDatabasePath(cfg.IPCityDatabasePath),
	)
	if err != nil {
		log.Printf("networking: failed to initialize ip intelligence client; ip routes disabled: %v", err)
	}

	vpnSvc := vpn.NewService(ipiClient)
	vpn.RegisterRoutes(r, vpnSvc)

	asnSvc := asn.NewService(ipiClient)
	asn.RegisterRoutes(r, asnSvc)

	infoSvc := info.NewService(ipiClient)
	info.RegisterRoutes(r, infoSvc)
}
