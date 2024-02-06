package routes

import (
	"github.com/divinitymn/aion-backend/internal/handlers"
  "github.com/divinitymn/aion-backend/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterPostRoutes(e *echo.Group) {
	posts := e.Group("/posts")

	posts.GET("", handlers.PostGetAll)
	posts.POST("", handlers.PostCreate, middlewares.VerifyToken)

	posts.GET("/:id", handlers.PostGetOne)
	posts.PUT("/:id", handlers.PostUpdate)
	posts.DELETE("/:id", handlers.PostDelete)
}
