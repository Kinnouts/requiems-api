package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"requiems-api/platform/config"
	"requiems-api/services/ai"
	"requiems-api/services/convert"
	"requiems-api/services/email"
	"requiems-api/services/entertainment"
	"requiems-api/services/finance"
	"requiems-api/services/fitness"
	"requiems-api/services/misc"
	"requiems-api/services/places"
	"requiems-api/services/tech"
	"requiems-api/services/text"
)

func registerV1Routes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client, cfg config.Config) {
	convertRouter := chi.NewRouter()
	convert.RegisterRoutes(convertRouter)
	r.Mount("/convert", convertRouter)

	textRouter := chi.NewRouter()
	text.RegisterRoutes(textRouter, pool)
	r.Mount("/text", textRouter)

	aiRouter := chi.NewRouter()
	ai.RegisterRoutes(aiRouter)
	r.Mount("/ai", aiRouter)

	emailRouter := chi.NewRouter()
	email.RegisterRoutes(emailRouter)
	r.Mount("/email", emailRouter)

	entertainmentRouter := chi.NewRouter()
	entertainment.RegisterRoutes(entertainmentRouter)
	r.Mount("/entertainment", entertainmentRouter)

	miscRouter := chi.NewRouter()
	misc.RegisterRoutes(ctx, miscRouter, pool, rdb)
	r.Mount("/misc", miscRouter)

	placesRouter := chi.NewRouter()
	places.RegisterRoutes(placesRouter)
	r.Mount("/places", placesRouter)

	techRouter := chi.NewRouter()
	tech.RegisterRoutes(techRouter, cfg)
	r.Mount("/tech", techRouter)

	financeRouter := chi.NewRouter()
	finance.RegisterRoutes(financeRouter, pool, rdb)
	r.Mount("/finance", financeRouter)

	fitnessRouter := chi.NewRouter()
	fitness.RegisterRoutes(fitnessRouter, pool)
	r.Mount("/fitness", fitnessRouter)
}
