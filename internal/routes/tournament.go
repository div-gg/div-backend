package routes

import (
	"github.com/divinitymn/div-backend/internal/handlers"
	"github.com/divinitymn/div-backend/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterTournamentRoutes(e *echo.Group) {
  tournament := e.Group("/tournament")

  tournament.GET("", handlers.TournamentGetAll)
  tournament.GET("", handlers.TournamentCreate, middlewares.VerifyToken)

  tournament.GET("/:id", handlers.TournamentGetByID)
  tournament.PUT("/:id", handlers.TournamentUpdateByID, middlewares.VerifyToken)
  tournament.DELETE("/:id", handlers.TournamentDeleteByID, middlewares.VerifyToken)
}
