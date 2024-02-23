package middlewares

import (
	"net/http"

	"github.com/divinitymn/aion-backend/internal/models"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}

		return nil
	}
}
