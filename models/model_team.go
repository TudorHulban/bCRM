package models

import (
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

// Team Structure holding team fields. Every application user belongs to a team.
type TeamFormData struct {
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
func NewTeam(c echo.Context, db *pg.DB, f TeamFormData, noValidation bool) (*Team, error) {
	// validate data
	if !noValidation {
		errValid := isValidStruct(f, c.Logger())
		if errValid != nil {
			return nil, errValid
		}
	}

	// check db connection. debug level = 1
	if c.Logger().Level() == 1 {
		errQuery := commons.CheckPgDB(c.Logger(), db)
		if errQuery != nil {
			return nil, errQuery
		}
	}
	c.Logger().Debugf("database is responding.")

	return &Team{
		TeamFormData: f,
		tools: tools{
			log: c.Logger(),
			db:  db,
		},
	}, nil
}

func (t *Team) Insert() error {
	t.log.Debugf("team data to insert: %v", t.TeamFormData)

	if errInsert := t.db.Insert(&t.TeamFormData); errInsert != nil {
		return errInsert
	}

	return nil
}
