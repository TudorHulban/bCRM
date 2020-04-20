package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

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
