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

  "github.com/labstack/echo/v4"
)


