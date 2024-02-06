package middlewares

import (
	"strings"

	"github.com/divinitymn/aion-backend/internal/utils"

	"github.com/labstack/echo/v4"
)

func VerifyToken(key string, c echo.Context) (bool, error) {
  header := c.Request().Header.Get("Authorization")
  splitheader := strings.Split(header, " ")
  token := splitheader[1]

  _ = utils.ParseToken(token)

  return true, nil
}
