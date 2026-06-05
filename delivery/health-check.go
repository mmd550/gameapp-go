package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func HealthCheckHandler(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
