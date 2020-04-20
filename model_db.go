package main

import (
	"github.com/go-pg/pg"
)

// IAccount Interface provides decoupling of persistence solution for account related operations.
type IStore interface {
	Open(DBConnInfo) (*pg.DB, error)
	Close(*pg.DB) error
	CreateUser(User) error
	GetUserByUserCode(string) (User, error)
}

type PgStore struct{}

func (s *PgStore) Open(info DBConnInfo) (*pg.DB, error) {
	return pg.Connect(&pg.Options{
		Addr:     info.Socket,
		User:     info.User,
		Password: info.Pass,
		Database: info.DB,
	}), nil
}

func (s *PgStore) Close(db *pg.DB) error {

	return nil
}
