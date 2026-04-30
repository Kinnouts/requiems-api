package technology

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"requiems-api/services/technology/barcode"
	"requiems-api/services/technology/base64"
	"requiems-api/services/technology/color"
	"requiems-api/services/technology/counter"
	convformat "requiems-api/services/technology/format"
	"requiems-api/services/technology/markdown"
	"requiems-api/services/technology/numbase"
	"requiems-api/services/technology/password"
	"requiems-api/services/technology/qr"
	randomuser "requiems-api/services/technology/random_user"
	"requiems-api/services/technology/units"
	"requiems-api/services/technology/useragent"
)

// RegisterRoutes wires all technology sub-services onto the given router.
func RegisterRoutes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client) {
	markdownSvc := markdown.NewService()
	markdown.RegisterRoutes(r, markdownSvc)

	base64Svc := base64.NewService()
	base64.RegisterRoutes(r, base64Svc)

	numbaseSvc := numbase.NewService()
	numbase.RegisterRoutes(r, numbaseSvc)

	formatSvc := convformat.NewService()
	convformat.RegisterRoutes(r, formatSvc)

	colorSvc := color.NewService()
	color.RegisterRoutes(r, colorSvc)

	unitsSvc := units.NewService()
	units.RegisterRoutes(r, unitsSvc)

	counterRepo := counter.NewRepository(pool)
	counterSvc := counter.NewService(rdb, counterRepo)
	go counter.StartSyncWorker(ctx, rdb, counterRepo)
	counter.RegisterRoutes(r, counterSvc)

	randomUserSvc := randomuser.NewService()
	randomuser.RegisterRoutes(r, randomUserSvc)

	barcodeSvc := barcode.NewService()
	barcode.RegisterRoutes(r, barcodeSvc)

	passwordSvc := password.NewService()
	password.RegisterRoutes(r, passwordSvc)

	qrSvc := qr.NewService()
	qr.RegisterRoutes(r, qrSvc)

	uaSvc := useragent.NewService()
	useragent.RegisterRoutes(r, uaSvc)
}
