package models

import (
	"testing"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func Test_Team_Insert(t *testing.T) {
	dbConn := pg.Connect(&pg.Options{
		Addr:     commons.DBSocket,
		User:     commons.DBUser,
		Password: commons.DBPass,
		Database: commons.DBName,
	})
	defer dbConn.Close()

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	ectx := e.NewContext(nil, nil)

	f := TeamFormData{CODE: "BLUE" + commons.UXSecs(), Name: "Blue Team" + commons.UXSecs(), Description: "Blue Team"}
	team, errNew := NewTeam(ectx, dbConn, f, false)
	if assert.NoError(t, errNew) {
		assert.NoError(t, team.Insert())
	}
}
