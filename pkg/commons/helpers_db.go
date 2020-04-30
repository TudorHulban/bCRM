package commons

import (
	"errors"

	"github.com/labstack/echo"
)

func CheckPgDB(log echo.Logger) error {
	log.Debugf("checking DB connectivity")

	if DB() == nil {
		return errors.New("db pointer is nil")
	}
	_, errQuery := DB().Exec("SELECT 1")
	return errQuery
}
