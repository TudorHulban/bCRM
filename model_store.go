package main

import (
	"github.com/go-pg/pg"
)

type PgStore struct {
	TheDB *pg.DB
}

func (s *PgStore) Open(info DBConnInfo) (*pg.DB, error) {
	if s.TheDB == nil {
		s.TheDB = pg.Connect(&pg.Options{
			Addr:     info.Socket,
			User:     info.User,
			Password: info.Pass,
			Database: info.DB,
		})
	}
	return s.TheDB, nil
}

func (s *PgStore) Close() error {
	s.TheDB.Close()
	return nil
}
