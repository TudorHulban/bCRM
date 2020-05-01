package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/TudorHulban/bCRM/models"
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
)

// NewTeam Handler to create a new team. To be used by user management roles.
// Information needed for creating a team is:
// Name, description, code
// RAW testing: curl -d "teamid=1&code=1234&pass=abcd" -X POST http://localhost:8001/newuser
func NewTeam(c echo.Context) error {
	var e httpError
	c.Logger().Debug("New Team")

	if len(c.FormValue(commons.NewTeamFormCODE)) == 0 {
		c.Logger().Debug("received ", commons.NewTeamFormCODE, " as: ", c.FormValue(commons.NewTeamFormCODE))
		e.TheError = commons.NewTeamFormCODE + " information is not valid"
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
	c.Logger().Debugf("user model instance created: %v", user)

	if errInsert := user.Insert(context.TODO(), commons.CTXTimeOutSecs); errInsert != nil {
		c.Logger().Debug("errInsert:", errInsert)
		e.TheError = errInsert.Error()
		c.JSON(http.StatusInternalServerError, e)
		return errInsert
	}

	result := struct {
		ID int64 `json:ID`
	}{
		ID: user.ID,
	}
	return c.JSON(http.StatusCreated, result)
}
