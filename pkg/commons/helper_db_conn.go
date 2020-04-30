package commons

import (
	"github.com/go-pg/pg/v9"
)

var dbConn *pg.DB

func DB() *pg.DB {
	if dbConn == nil {
		dbConn = pg.Connect(&pg.Options{
			Addr:     DBSocket,
			User:     DBUser,
			Password: DBPass,
			Database: DBName,
		})
	}
	return dbConn
}
