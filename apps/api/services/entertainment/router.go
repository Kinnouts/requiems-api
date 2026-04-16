package entertainment

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"requiems-api/services/entertainment/advice"
	"requiems-api/services/entertainment/chucknorris"
	"requiems-api/services/entertainment/emoji"
	"requiems-api/services/entertainment/facts"
	"requiems-api/services/entertainment/horoscope"
	"requiems-api/services/entertainment/jokes"
	"requiems-api/services/entertainment/quotes"
	"requiems-api/services/entertainment/sudoku"
	"requiems-api/services/entertainment/trivia"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool) {
	adviceSvc := advice.NewService(pool)
	advice.RegisterRoutes(r, adviceSvc)

	quotesSvc := quotes.NewService(pool)
	quotes.RegisterRoutes(r, quotesSvc)

	horoscopeSvc := horoscope.NewService()
	horoscope.RegisterRoutes(r, horoscopeSvc)

	sudokuSvc := sudoku.NewService()
	sudoku.RegisterRoutes(r, sudokuSvc)

	emojiSvc := emoji.NewService()
	emoji.RegisterRoutes(r, emojiSvc)

	factsSvc := facts.NewService()
	facts.RegisterRoutes(r, factsSvc)

	triviaSvc := trivia.NewService()
	trivia.RegisterRoutes(r, triviaSvc)

	chuckNorrisSvc := chucknorris.NewService()
	chucknorris.RegisterRoutes(r, chuckNorrisSvc)

	jokesSvc := jokes.NewService()
	jokes.RegisterRoutes(r, jokesSvc)
}
