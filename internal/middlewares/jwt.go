package middlewares

import (
  "time"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/divinitymn/aion-backend/internal/models"
	"github.com/divinitymn/aion-backend/internal/utils"
)

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
    token := strings.Split(c.Request().Header.Get("Authorization"), "Bearer ")
    claims := utils.ParseToken(token[1])
    valid := claims.VerifyExpiresAt(time.Now().Unix(), true)

    if valid {
      return next(c)
    } else {
      return c.JSON(http.StatusUnauthorized, models.Response{
        Status: http.StatusUnauthorized,
        Message: "Token expired",
      })
    }
	}
}

