package main

import (
	"github.com/go-pg/pg/v9"
)

func CheckPgDB(db *pg.DB) error {
	_, errQuery := db.Exec("SELECT 1")
	return errQuery
}
