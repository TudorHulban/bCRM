package main

import (
	"net/http"
	"testing"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func Test1CreateUser(t *testing.T) {
	dbConn := pg.Connect(&pg.Options{
		Addr:     commons.DBSocket,
		User:     commons.DBUser,
		Password: commons.DBPass,
		Database: commons.DBName,
	})
	defer dbConn.Close()

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointNewUser, NewUser)

	if assert.Nil(t, commons.CheckPgDB(e.Logger, dbConn), "TEST - Could not connect to DB.") {
		apitest.New().
			Handler(e).
			Method(http.MethodPost).
			URL(commons.EndpointNewUser).
			FormData("teamid", "1").
			FormData("code", "MARY").
			FormData("pass", "abcd").
			Expect(t).
			Status(http.StatusOK).
			Body(`{"ok":"OK"}`).
			End()
	}
}
