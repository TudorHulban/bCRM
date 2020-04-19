package interfaces

import (
	"github.com/TudorHulban/bCRM/structs"
	"github.com/go-pg/pg"
)

type DBConnInfo struct {
	Socket string
	User   string
	Pass   string
	DB     string
}

// when working with different RDBMS change.
type Database *pg.DB

// IAccount Interface provides decoupling of persistence solution for account related operations.
type IStore interface {
	Open(DBConnInfo) (Database, error)
	Close(Database) error
	CreateUser(structs.User) error
	GetUserByUserCode(string) (structs.User, error)
}
