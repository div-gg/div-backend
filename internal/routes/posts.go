package routes

import (
  "github.com/divinitymn/aion-backend/internal/handlers"

  "github.com/labstack/echo/v4"
)

func RegisterPostRoutes(e *echo.Group) {
  posts := e.Group("/posts")

  posts.GET("/", handlers.PostGetAll)
  posts.POST("/", handlers.PostCreate)

  posts.GET("/:id", handlers.PostGetSingle)
  posts.PUT("/:id", handlers.PostUpdate)
  posts.DELETE("/:id", handlers.PostDelete)
}
