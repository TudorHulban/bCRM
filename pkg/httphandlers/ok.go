package httphandlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func HandlerOK(c echo.Context) error {
	log.Info("OK")
	return c.JSON(http.StatusOK, echo.Map{"ok": "OK"})
}
