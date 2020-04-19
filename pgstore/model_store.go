package pgstore

import (
	"github.com/TudorHulban/bCRM/interfaces"
	"github.com/go-pg/pg"
)

// PgStore Satisfies IStore interface.
type PgStore struct{}

func (s *PgStore) Open(info interfaces.DBConnInfo) (interfaces.Database, error) {
	return pg.Connect(&pg.Options{
		Addr:     info.Socket,
		User:     info.User,
		Password: info.Pass,
		Database: info.DB,
	}), nil
}

func (s *PgStore) Close(db interfaces.Database) error {
	db.Close
}
