package main

import (
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
)

func TestUserAdd(t *testing.T) {
	dbconnInfo := DBConnInfo{
		Socket: DBSocket,
		User:   DBUser,
		Pass:   DBPass,
		DB:     DBName,
	}
	dbConn := pg.Connect(&pg.Options{
		Addr:     dbconnInfo.Socket,
		User:     dbconnInfo.User,
		Password: dbconnInfo.Pass,
		Database: dbconnInfo.DB,
	})
	defer dbConn.Close()

	if assert.NoError(t, CheckPgDB(dbConn), "no connection to DB") {
		u := User{ID: 1, LoginCODE: "xxx", SecurityGroup: 1}
		errAdd := CreateUser(&u)

		if errAdd != nil {
			t.Error("TestUserAdd:", errAdd)
		}
	}

}
