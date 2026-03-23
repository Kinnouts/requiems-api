package entertainment

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/entertainment/chucknorris"
	"requiems-api/services/entertainment/emoji"
	"requiems-api/services/entertainment/facts"
	"requiems-api/services/entertainment/horoscope"
	"requiems-api/services/entertainment/jokes"
	"requiems-api/services/entertainment/sudoku"
	"requiems-api/services/entertainment/trivia"
)

func RegisterRoutes(r chi.Router) {
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
