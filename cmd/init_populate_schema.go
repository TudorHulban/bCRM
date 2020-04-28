package main

import (
	"context"

	"github.com/TudorHulban/bCRM/pkg/commons"

	"github.com/TudorHulban/bCRM/models"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

func initTeams(ctx context.Context, c echo.Context, db *pg.DB) error {
	teams := []models.TeamFormData{}
	teams = append(teams, models.TeamFormData{CODE: "APP", Name: "APP Admins", Description: "APP Admins Team"})
	teams = append(teams, models.TeamFormData{CODE: "BLUE", Name: "Blue Team", Description: "Blue Team"})
	teams = append(teams, models.TeamFormData{CODE: "YELLOW", Name: "Yellow Team", Description: "Yellow Team"})

	for _, v := range teams {
		t, errCo := models.NewTeam(c, db, v, false)
		if errCo != nil {
			return errCo
		}
		errInsert := t.Insert(ctx, commons.CTXTimeOutSecs)
		if errInsert != nil {
			return errInsert
		}
	}
	return nil
}

func initGroups(ctx context.Context, c echo.Context, db *pg.DB) error {
	groups := []models.GroupFormData{}
	groups = append(groups, models.GroupFormData{CODE: "APP", Name: "APP Admins Group", Description: "Group for app admins."})
	groups = append(groups, models.GroupFormData{CODE: "G1", Name: "DBA Group", Description: "Group for database administrators."})
	groups = append(groups, models.GroupFormData{CODE: "G2", Name: "Support Group", Description: "Group for support engineers."})

	for _, v := range groups {
		g, errCo := models.NewGroup(c, db, v, false)
		if errCo != nil {
			return errCo
		}
		errInsert := g.Insert(ctx, commons.CTXTimeOutSecs)
		if errInsert != nil {
			return errInsert
		}
	}
	return nil
}

func initUsers(ctx context.Context, c echo.Context, db *pg.DB) error {
	users := []models.UserFormData{}
	users = append(groups, models.GroupFormData{CODE: "ADMIN", Name: "APP Admin", Description: "Account for application admin."})
	users = append(groups, models.GroupFormData{CODE: "USER1", Name: "John SMith", Description: "Test account for user w/o rights."})

	for _, v := range groups {
		u, errCo := models.NewUser(c, db, v, false)
		if errCo != nil {
			return errCo
		}
		errInsert := u.Insert(ctx, commons.CTXTimeOutSecs)
		if errInsert != nil {
			return errInsert
		}
	}
	return nil
}
