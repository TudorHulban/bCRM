package models

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

// tools Structure providing validation, logging and database connection to models.
type tools struct {
	log echo.Logger
	db  *pg.DB
}
