package commons

import (
	"errors"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

func CheckPgDB(log echo.Logger, db *pg.DB) error {
	log.Debugf("checking DB connectivity")

	if db == nil {
		return errors.New("db pointer is nil")
	}
	_, errQuery := db.Exec("SELECT 1")
	return errQuery
}
