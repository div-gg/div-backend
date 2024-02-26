package routes

import (
	"github.com/divinitymn/div-backend/internal/handlers"
	"github.com/divinitymn/div-backend/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterPostRoutes(e *echo.Group) {
	posts := e.Group("/posts")

	posts.GET("", handlers.PostGetAll)
	posts.POST("", handlers.PostCreate, middlewares.VerifyToken)

	posts.GET("/:id", handlers.PostGetByID)
	posts.PUT("/:id", handlers.PostUpdateByID, middlewares.VerifyToken)
	posts.DELETE("/:id", handlers.PostDeleteByID, middlewares.VerifyToken)
}
