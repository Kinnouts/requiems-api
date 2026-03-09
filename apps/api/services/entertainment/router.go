package entertainment

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/entertainment/emoji"
	"requiems-api/services/entertainment/horoscope"
	"requiems-api/services/entertainment/sudoku"
)

func RegisterRoutes(r chi.Router) {
	horoscopeSvc := horoscope.NewService()
	horoscope.RegisterRoutes(r, horoscopeSvc)

	sudokuSvc := sudoku.NewService()
	sudoku.RegisterRoutes(r, sudokuSvc)
  
	emojiSvc := emoji.NewService()
	emoji.RegisterRoutes(r, emojiSvc)
}
