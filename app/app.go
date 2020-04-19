package app

import (
	"github.com/go-pg/pg"
)

type database *pg.DB

type App struct {
	DB database
}

func NewApp() *App {
	return &App{}
}
