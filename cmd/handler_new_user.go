package main

import (
	"net/http"

	"github.com/TudorHulban/bCRM/models"
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo"
)

// NewUser Handler to create a new user. To be used by user management roles.
// Information needed for creating a user is:
// Name, UserCode, Password
func NewUser(c echo.Context) error {
	var e httpError
	c.Logger().Debug("New User")

	if len(c.FormValue(commons.NewUserFormName)) == 0 {
		msg := commons.NewUserFormName + " information is not valid"
		c.Logger().Debugf("%v", msg)
		e.TheError = msg
		return c.JSON(http.StatusNotAcceptable, e)
	}

	if len(c.FormValue(commons.NewUserFormUserCode)) == 0 {
		e.TheError = commons.NewUserFormUserCode + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	if len(c.FormValue(commons.NewUserFormPass)) == 0 {
		e.TheError = commons.NewUserFormPass + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	var co models.Contact
	co.FirstName = c.FormValue(commons.NewUserFormName)
	co.CompanyEmail = c.FormValue(commons.NewUserFormEmail)

	c.Logger().Debug("Contact:", co)

	var u models.UserFormData
	u.LoginCODE = c.FormValue(commons.NewUserFormUserCode)
	u.LoginPWD = c.FormValue(commons.NewUserFormPass)
	u.SecurityGroup = commons.SecuGrpUser

	c.Logger().Debug("User:", u)

	user := models.NewUser(c, dbConn, u)
	errInsert := user.Insert()
	if errInsert != nil {
		c.Logger().Debug("errAdd:", errInsert)
		e.TheError = errInsert.Error()
		c.JSON(http.StatusInternalServerError, e)
		return errInsert
	}
	return c.JSON(http.StatusOK, user.ID)
}
