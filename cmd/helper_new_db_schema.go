package main

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
)

// newSchema Creates tables in PostgreSQL with the help of the ORM and based on passed models.
func newSchema(db *pg.DB, models []interface{}) error {
	createTable4Model := func(model interface{}) error {
		return db.CreateTable(model, &orm.CreateTableOptions{Temp: false, IfNotExists: true, FKConstraints: true})
	}

	for _, model := range models {
		if errCreateTable := createTable4Model(model); errCreateTable != nil {
			return errCreateTable
		}
	}
	return nil
}

func populateSchema(ctx context.Context, c echo.Context, db *pg.DB) error {
	if errTeams := initTeams(ctx, c, db); errTeams != nil {
		return errTeams
	}
	if errGroups := initGroups(ctx, c, db); errGroups != nil {
		return errGroups
	}
	if errUsers := initUsers(ctx, c, db); errUsers != nil {
		return errUsers
	}
	return nil
}
