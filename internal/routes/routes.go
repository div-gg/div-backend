package routes

import (
	"net/http"

	"github.com/divinitymn/div-backend/internal/models"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.Response{
			Status:  http.StatusOK,
			Message: "OK",
		})
	})

	RegisterAuthRoutes(v1)
	RegisterPostRoutes(v1)
	RegisterUserRoutes(v1)
	RegisterTournamentRoutes(v1)
}
