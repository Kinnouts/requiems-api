package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"requiems-api/internal/email"
	"requiems-api/internal/entertainment"
	"requiems-api/internal/misc"
	"requiems-api/internal/places"
	"requiems-api/internal/tech"
	"requiems-api/internal/text"
)

func registerV1Routes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client) {
	textRouter := chi.NewRouter()
	text.RegisterRoutes(textRouter, pool)
	r.Mount("/text", textRouter)

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
	tech.RegisterRoutes(techRouter)
	r.Mount("/tech", techRouter)
}
