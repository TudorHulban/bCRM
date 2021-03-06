package models

import (
	"context"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

func TableExists(ctx context.Context, model interface{}, table string, timeoutSecs int, log echo.Logger) (bool, error) {
	dml := "select exists(select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + table + "'" + ")"

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	log.Debug("DML:", dml)

	var exists bool
	_, errExists := commons.DB().QueryOneContext(ctx, pg.Scan(&exists), dml)
	if errExists != nil {
		log.Debug("errExists:", errExists)
		return false, errExists
	}
	log.Debug("Exists:", errExists)
	return exists, nil
}
