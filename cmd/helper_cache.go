package main

import (
	"github.com/TudorHulban/bCRM/pkg/cache"
	"github.com/labstack/echo"
)

var c *cache.Cache

// getSessionIDCache Helper to create cache. Call it only once with logger.
func getSessionIDCache(log ...echo.Logger) *cache.Cache {
	if c.DB == nil && len(log) > 0 {
		var errorCache error
		c, errorCache = cache.New(log[0])
	}
	return c
}
