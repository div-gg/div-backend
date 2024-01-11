package routes

import (
  "github.com/divinitymn/aion-backend/internal/handlers"

  "github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Group) {
  auth := e.Group("/auth")

  auth.POST("/register", handlers.RegisterHandler)
  auth.POST("/login", handlers.LoginHandler)
}
