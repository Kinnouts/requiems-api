package app

import (
	"context"
	"strings"

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

// serviceEnabled reports whether a named service should be mounted.
// When EnabledServices is empty (the default), all services are mounted —
// preserving the existing shared-deployment behaviour unchanged.
// For private deployments, set ENABLED_SERVICES="email,text,tech" to mount
// only the services the tenant purchased.
func serviceEnabled(cfg config.Config, key string) bool {
	if strings.TrimSpace(cfg.EnabledServices) == "" {
		return true
	}
	for s := range strings.SplitSeq(cfg.EnabledServices, ",") {
		if strings.TrimSpace(s) == key {
			return true
		}
	}
	return false
}

func registerV1Routes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client, cfg config.Config) {
	if serviceEnabled(cfg, "convert") {
		convertRouter := chi.NewRouter()
		convert.RegisterRoutes(convertRouter)
		r.Mount("/convert", convertRouter)
	}

	if serviceEnabled(cfg, "text") {
		textRouter := chi.NewRouter()
		text.RegisterRoutes(textRouter, pool)
		r.Mount("/text", textRouter)
	}

	if serviceEnabled(cfg, "ai") {
		aiRouter := chi.NewRouter()
		ai.RegisterRoutes(aiRouter)
		r.Mount("/ai", aiRouter)
	}

	if serviceEnabled(cfg, "email") {
		emailRouter := chi.NewRouter()
		email.RegisterRoutes(emailRouter)
		r.Mount("/email", emailRouter)
	}

	if serviceEnabled(cfg, "entertainment") {
		entertainmentRouter := chi.NewRouter()
		entertainment.RegisterRoutes(entertainmentRouter)
		r.Mount("/entertainment", entertainmentRouter)
	}

	if serviceEnabled(cfg, "misc") {
		miscRouter := chi.NewRouter()
		misc.RegisterRoutes(ctx, miscRouter, pool, rdb)
		r.Mount("/misc", miscRouter)
	}

	if serviceEnabled(cfg, "places") {
		placesRouter := chi.NewRouter()
		places.RegisterRoutes(placesRouter, cfg, rdb)
		r.Mount("/places", placesRouter)
	}

	if serviceEnabled(cfg, "tech") {
		techRouter := chi.NewRouter()
		tech.RegisterRoutes(techRouter, cfg)
		r.Mount("/tech", techRouter)
	}

	if serviceEnabled(cfg, "finance") {
		financeRouter := chi.NewRouter()
		finance.RegisterRoutes(financeRouter, pool, rdb)
		r.Mount("/finance", financeRouter)
	}

	if serviceEnabled(cfg, "fitness") {
		fitnessRouter := chi.NewRouter()
		fitness.RegisterRoutes(fitnessRouter, pool)
		r.Mount("/fitness", fitnessRouter)
	}
}
