package routes

import (
  "github.com/divinitymn/div-backend/internal/handlers"
  "github.com/divinitymn/div-backend/internal/middlewares"

  "github.com/labstack/echo/v4"
)

func RegisterBlogRoutes(e *echo.Group) {
  blog := e.Group("/blogs")

  blog.GET("", handlers.BlogGetAll)
  blog.POST("", handlers.BlogCreate, middlewares.VerifyToken)

  blog.GET("/:id", handlers.BlogGetByID)
  blog.PUT("/:id", handlers.BlogUpdateByID, middlewares.VerifyToken)
  blog.DELETE("/:id", handlers.BlogDeleteByID, middlewares.VerifyToken)
}
