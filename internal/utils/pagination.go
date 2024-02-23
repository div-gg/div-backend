package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetPaginationValues(c echo.Context) (int64, int64) {
	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	page = page - 1
	if err != nil || page < 0 {
		page = 0
	}

	limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err != nil {
		limit = 10
	}

	return page, limit
}
