package middlewares

import (
	"net/http"
	"strings"
	"time"
  "log"

	"github.com/labstack/echo/v4"

	"github.com/divinitymn/div-backend/internal/models"
	"github.com/divinitymn/div-backend/internal/utils"
)

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := strings.Split(c.Request().Header.Get("Authorization"), "Bearer ")

    if len(token) != 2 {
      return c.JSON(http.StatusUnauthorized, models.Response{
        Status:  http.StatusUnauthorized,
        Message: "Invalid token",
      })
    }

		claims := utils.ParseToken(token[1])
		valid := claims.VerifyExpiresAt(time.Now().Unix(), true)
    log.Println(claims)

    c.Set("userId", claims["id"])

		if valid {
			return next(c)
		} else {
			return c.JSON(http.StatusUnauthorized, models.Response{
				Status:  http.StatusUnauthorized,
				Message: "Token expired",
			})
		}
	}
}
