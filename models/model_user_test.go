package models

import (
	"testing"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestUserAdd(t *testing.T) {
	dbConn := pg.Connect(&pg.Options{
		Addr:     commons.DBSocket,
		User:     commons.DBUser,
		Password: commons.DBPass,
		Database: commons.DBName,
	})
	defer dbConn.Close()

	e := echo.New()
	ectx := e.NewContext(nil, nil)

	f := UserFormData{TeamID: 1, LoginCODE: "xxx", SecurityGroup: 1, LoginPWD: "abcd"}
	user := NewUser(ectx, dbConn, f)

	err := user.Insert()
	assert.NoError(t, err)

}
