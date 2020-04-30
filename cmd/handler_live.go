package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Live HTTP Endpoint for verifying app is operational.
func Live(c echo.Context) error {
	c.Logger().Debug("Live")
	return c.JSON(http.StatusOK, echo.Map{"ok": "OK"})
}
