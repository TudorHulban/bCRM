package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

func TableExists(ctx context.Context, db *pg.DB, table string, timeoutSecs int, log echo.Logger) error {
	dml := "SELECT exists (select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + table + "'" + ")"
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	result, err := db.ExecOneContext(ctx, dml)
	if err != nil {
		return err
	}
	log.Debug("result:", result)
	return nil
}
