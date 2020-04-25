package models

import (
	"testing"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestUserInsert(t *testing.T) {
	dbConn := pg.Connect(&pg.Options{
		Addr:     commons.DBSocket,
		User:     commons.DBUser,
		Password: commons.DBPass,
		Database: commons.DBName,
	})
	defer dbConn.Close()

	e := echo.New()
	ectx := e.NewContext(nil, nil)

	f := UserFormData{TeamID: 1, LoginCODE: "xxx3", SecurityGroup: 1, LoginPWD: "abcd"}
	user, errNew := NewUser(ectx, dbConn, f, false)
	assert.NoError(t, errNew)

	assert.NoError(t, user.Insert())
}

func TestUserSelectByID(t *testing.T) {
	dbConn := pg.Connect(&pg.Options{
		Addr:     commons.DBSocket,
		User:     commons.DBUser,
		Password: commons.DBPass,
		Database: commons.DBName,
	})
	defer dbConn.Close()

	e := echo.New()
	ectx := e.NewContext(nil, nil)

	f := UserFormData{}
	user, errNew := NewUser(ectx, dbConn, f, true)
	assert.NoError(t, errNew)

	_, err := user.GetbyID(1)
	assert.NoError(t, err)
}
