package main

import (
	"net/http"
	"strconv"

	"github.com/TudorHulban/bCRM/models"
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo"
)

// NewUser Handler to create a new user. To be used by user management roles.
// Information needed for creating a user is:
// Name, UserCode, Password
// RAW testing: curl -d "teamid=1&code=1234&pass=abcd" -X POST http://localhost:8001/newuser
func NewUser(c echo.Context) error {
	var e httpError
	c.Logger().Debug("New User")

	if len(c.FormValue(commons.NewUserFormTeamID)) == 0 {
		c.Logger().Debug("received ", commons.NewUserFormTeamID, " as: ", c.FormValue(commons.NewUserFormTeamID))
		e.TheError = commons.NewUserFormTeamID + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}
	if len(c.FormValue(commons.NewUserFormUserCode)) == 0 {
		c.Logger().Debug("received ", commons.NewUserFormUserCode, " as: ", c.FormValue(commons.NewUserFormUserCode))
		e.TheError = commons.NewUserFormUserCode + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}
	if len(c.FormValue(commons.NewUserFormPass)) == 0 {
		c.Logger().Debug("received ", commons.NewUserFormPass, " as: ", c.FormValue(commons.NewUserFormPass))
		e.TheError = commons.NewUserFormPass + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	var u models.UserFormData
	id, errConv := strconv.Atoi(c.FormValue(commons.NewUserFormTeamID))
	if errConv != nil {
		e.TheError = commons.NewUserFormTeamID + " could not convert"
		return c.JSON(http.StatusNotAcceptable, e)
	}
	u.TeamID = id
	u.LoginCODE = c.FormValue(commons.NewUserFormUserCode)
	u.LoginPWD = c.FormValue(commons.NewUserFormPass)
	u.SecurityGroup = commons.SecuGrpUser

	c.Logger().Debug("User:", u)

	// check db connection for debug level
	if c.Logger().Level() == 1 {
		errQuery := commons.CheckPgDB(c.Logger())
		c.Logger().Debug("errQuery:", errQuery)
		if errQuery != nil {
			return errQuery
		}
		c.Logger().Debugf("database is responding.")
	}

	user, errCo := models.NewUser(c, u, false)
	if errCo != nil {
		c.Logger().Debug("errCo:", errCo)
		e.TheError = errCo.Error()
		c.JSON(http.StatusServiceUnavailable, e)
		return errCo
	}

	errInsert := user.Insert(ctx, commons.CTXTimeOutSecs)
	if errInsert != nil {
		c.Logger().Debug("errInsert:", errInsert)
		e.TheError = errInsert.Error()
		c.JSON(http.StatusInternalServerError, e)
		return errInsert
	}
	return c.JSON(http.StatusOK, user.ID)
}
