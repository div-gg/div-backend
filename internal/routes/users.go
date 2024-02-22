package routes

import (
	"github.com/divinitymn/aion-backend/internal/handlers"
  // "github.com/divinitymn/aion-backend/internal/middlewares"

  "github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Group) {
  users := e.Group("/posts")

  users.GET("/:id", handlers.UserGetByID)
}
