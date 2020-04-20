package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// NewUser Creates a new user. To be used by user management roles.
// Information needed for creating a user is:
// Name, UserCode, Password
func NewUser(c echo.Context) error {
	var e httpError

	if len(c.FormValue(NewUserFormName)) == 0 {
		msg := NewUserFormName + " information is not valid"
		c.Logger().Debugf("%v", msg)
		e.TheError = msg
		return c.JSON(http.StatusNotAcceptable, e)
	}

	if len(c.FormValue(NewUserFormUserCode)) == 0 {
		e.TheError = NewUserFormUserCode + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	if len(c.FormValue(NewUserFormPass)) == 0 {
		e.TheError = NewUserFormPass + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	var co Contact
	co.FirstName = c.FormValue(NewUserFormName)
	co.CompanyEmail = c.FormValue(NewUserFormEmail)

	var u User
	u.LoginCODE = c.FormValue(NewUserFormUserCode)
	u.LoginPWD = c.FormValue(NewUserFormPass)
	u.ContactInfo = append(u.ContactInfo, &co)
	u.SecurityGroup = SecuGrpUser

	errAdd := store.CreateUser(&u)
	if errAdd != nil {
		e.TheError = errAdd.Error()
		return c.JSON(http.StatusInternalServerError, e)
	}
	return c.JSON(http.StatusOK, u.ID)
}