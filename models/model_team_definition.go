package models

import (
	"context"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Team Structure holding team fields. Every application user belongs to a team.
type TeamFormData struct {
	tableName       struct{} `pg:"teams"`
	ID              int64
	CODE            string `validate:"required" pg:",notnull,unique"`
	Name            string `validate:"required" pg:",notnull,unique"`
	Description     string `validate:"required" pg:",notnull"`
	AssignedTickets int    `pg:"numbtickets"`
	ManagerID       int64  `pg:"managerid"`
}

type Team struct {
	TeamFormData
	tools
}

// NewTeam Constructor for when interacting with the model.
func NewTeam(c echo.Context, f TeamFormData, noValidation bool) (*Team, error) {
	// validate data
	if !noValidation {
		if errValid := isValidStruct(f, c.Logger()); errValid != nil {
			return nil, errValid
		}
	}

	// check db connection. debug level = 1
	if c.Logger().Level() == 1 {
		if errQuery := commons.CheckPgDB(c.Logger()); errQuery != nil {
			return nil, errQuery
		}
		c.Logger().Debugf("database is responding.")
	}

	result := Team{
		TeamFormData: f,
		tools: tools{
			log: c.Logger(),
			db:  commons.DB(),
		},
	}
	result.tools.log.SetLevel(log.DEBUG)
	return &result, nil
}

func (t *Team) Insert(ctx context.Context, timeoutSecs int) error {
	t.log.Debugf("team data to insert: %v", t.TeamFormData)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	if errInsert := t.db.WithContext(ctx).Insert(&t.TeamFormData); errInsert != nil {
		return errInsert
	}
	return nil
}
