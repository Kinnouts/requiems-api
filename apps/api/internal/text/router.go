package text

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"requiems-api/internal/text/advice"
	"requiems-api/internal/text/lorem"
	"requiems-api/internal/text/profanity"
	"requiems-api/internal/text/quotes"
	"requiems-api/internal/text/words"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool) {
	adviceSvc := advice.NewService(pool)
	advice.RegisterRoutes(r, adviceSvc)

	quotesSvc := quotes.NewService(pool)
	quotes.RegisterRoutes(r, quotesSvc)

	wordsSvc := words.NewService(pool)
	words.RegisterRoutes(r, wordsSvc)

	loremSvc := lorem.NewService()
	lorem.RegisterRoutes(r, loremSvc)

	profanitySvc := profanity.NewService()
	profanity.RegisterRoutes(r, profanitySvc)
}
