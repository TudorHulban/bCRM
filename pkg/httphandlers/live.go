package httphandlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Live(c echo.Context) error {
	c.Logger().Debug("Live")
	return c.JSON(http.StatusOK, echo.Map{"ok": "OK"})
}
