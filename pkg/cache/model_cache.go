package cache

import (
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/labstack/echo/v4"
)

// Cache Structure implementing cache methods.
type Cache struct {
	DB  *badger.DB
	log echo.Logger
}

// New Constructor for in memory cache.
func New(log echo.Logger) *Cache {
	options := badger.DefaultOptions("").WithInMemory(true)
	result, errOpen := badger.Open(options)
	if errOpen != nil {
		return nil
	}

	return &Cache{
		log: log,
		DB:  result,
	}
}

// Close Method closes the cache.
func (c *Cache) Close() error {
	return c.DB.Close()
}

// SetTTL Method can be used for inserts and updates. Time To Live in seconds.
func (c *Cache) SetTTL(key, value string, ttlSeconds int) error {
	return c.DB.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry([]byte(key), []byte(value)).WithTTL(time.Second * time.Duration(ttlSeconds)))
	})
}

// Get Method fetches key from store.
func (c *Cache) Get(key string) ([]byte, error) {
	var result []byte

	errView := c.DB.View(func(txn *badger.Txn) error {
		item, errGet := txn.Get([]byte(key))
		if errGet != nil {
			return errGet
		}
		c.log.Debugf("size: %v, expires: %v", item.EstimatedSize(), item.ExpiresAt())

		errItem := item.Value(func(itemVals []byte) error {
			result = append(result, itemVals...)
			return nil
		})
		return errItem
	})
	return result, errView
}
