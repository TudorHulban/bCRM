package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// NewSchema Creates tables in PostgreSQL with the help of the ORM and based on passed models.
func NewSchema(db *pg.DB, models ...interface{}) error {
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
