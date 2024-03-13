package routes

import (
	"github.com/divinitymn/div-backend/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Group) {
	auth := e.Group("/auth")

	auth.POST("/login", handlers.LoginHandler)
  auth.POST("/register", handlers.RegisterHandler)
	auth.GET("/discord/callback", handlers.DiscordCallbackHandler)
}
