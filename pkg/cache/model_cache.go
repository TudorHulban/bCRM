package cache

import (
	badger "github.com/dgraph-io/badger/v2"
)

type Cache struct {
	db *badger.DB
}
