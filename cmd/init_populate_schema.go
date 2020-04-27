package main

import (
	"context"

	"github.com/TudorHulban/bCRM/models"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

func initTeams(ctx context.Context, c echo.Context, db *pg.DB) error {
	teams := []models.TeamFormData{}
	teams = append(teams, models.TeamFormData{CODE: "BLUE", Name: "Blue Team", Description: "Blue Team"})
	teams = append(teams, models.TeamFormData{CODE: "YELLOW", Name: "Yellow Team", Description: "Yellow Team"})

	for _, v := range teams {
		t, errCo := models.NewTeam(c, db, v, false)
		if errCo != nil {
			return errCo
		}
		errInsert := t.Insert()
		if errInsert != nil {
			return errInsert
		}
	}
	return nil
}
