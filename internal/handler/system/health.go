package system

import (
	"gosir/internal/common"

	"github.com/labstack/echo/v4"
)

func HealthCheck(c echo.Context) error {
	return common.Success(c, map[string]string{
		"status": "ok",
	})
}
