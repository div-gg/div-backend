package handlers

// posts.GET("/", handlers.PostGetAll)
// posts.POST("/", handlers.PostCreate)
//
// posts.GET("/:id", handlers.PostGetSingle)
// posts.PUT("/:id", handlers.PostUpdate)
// posts.DELETE("/:id", handlers.PostDelete)

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/divinitymn/aion-backend/internal/models"
	"github.com/divinitymn/aion-backend/internal/db"

	"github.com/labstack/echo/v4"
)

func PostGetAll(c echo.Context) error {
  ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
  defer cancel()

  return &echo.HTTPError{}
}

func PostGetSingle(c echo.Context) error {
  return &echo.HTTPError{}
}

func PostCreate(c echo.Context) error {
  return &echo.HTTPError{}
}

func PostUpdate(c echo.Context) error {
  return &echo.HTTPError{}
}

func PostDelete(c echo.Context) error {
  return &echo.HTTPError{}
}
