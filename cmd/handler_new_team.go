package main

import (
	"context"
	"net/http"

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
	if len(c.FormValue(commons.NewTeamFormDescription)) == 0 {
		c.Logger().Debug("received ", commons.NewTeamFormDescription, " as: ", c.FormValue(commons.NewTeamFormDescription))
		e.TheError = commons.NewTeamFormDescription + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}
	if len(c.FormValue(commons.NewTeamFormName)) == 0 {
		c.Logger().Debug("received ", commons.NewTeamFormName, " as: ", c.FormValue(commons.NewTeamFormName))
		e.TheError = commons.NewTeamFormName + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	var m models.TeamFormData
	m.Name = c.FormValue(commons.NewTeamFormName)
	m.Description = c.FormValue(commons.NewTeamFormDescription)
	m.CODE = c.FormValue(commons.NewTeamFormCODE)

	c.Logger().Debug("Team:", m)

	// check db connection for debug level
	if c.Logger().Level() == 1 {
		errQuery := commons.CheckPgDB(c.Logger())
		c.Logger().Debug("errQuery:", errQuery)
		if errQuery != nil {
			return errQuery
		}
		c.Logger().Debugf("database is responding.")
	}

	model, errCo := models.NewTeam(c, m, false)
	if errCo != nil {
		c.Logger().Debug("errCo:", errCo)
		e.TheError = errCo.Error()
		c.JSON(http.StatusServiceUnavailable, e)
		return errCo
	}
	c.Logger().Debugf("model instance created: %v", model)

	if errInsert := model.Insert(context.TODO(), commons.CTXTimeOutSecs); errInsert != nil {
		c.Logger().Debug("errInsert:", errInsert)
		e.TheError = errInsert.Error()
		c.JSON(http.StatusInternalServerError, e)
		return errInsert
	}

	result := struct {
		ID int64 `json:ID`
	}{
		ID: model.ID,
	}
	return c.JSON(http.StatusCreated, result)
}
