package routes

import (
	"github.com/divinitymn/div-backend/internal/handlers"
	"github.com/divinitymn/div-backend/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Group) {
	users := e.Group("/users")

	users.GET("/:id", handlers.UserGetByID)
  users.GET("/me", handlers.UserGetMe, middlewares.VerifyToken)
  users.PUT("/me", handlers.UserUpdateMe, middlewares.VerifyToken)
  // users.PUT("/:id", handlers.UserUpdateByID, middlewares.VerifyToken)
  // users.DELETE("/:id", handlers.UserDeleteByID, middlewares.VerifyToken)
}
