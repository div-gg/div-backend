package routes

import (
  "github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
  v1 := e.Group("/v1")
  RegisterAuthRoutes(v1)
}
