package models

import (
	"context"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

// GroupFormData Structure holding team fields. A team could belong to a group.
type GroupFormData struct {
	tableName       struct{} `pg:"groups"`
	ID              int64
	CODE            string `validate:"required" pg:",notnull,unique"`
	Name            string `validate:"required" pg:",notnull,unique"`
	Description     string `validate:"required" pg:",notnull"`
	AssignedTickets int    `pg:"numbtickets"`
	ManagerID       int64  `pg:"managerid"`
}

type Group struct {
	GroupFormData
	tools
}

// NewGroup Constructor for when interacting with the model.
func NewGroup(c echo.Context, db *pg.DB, f GroupFormData, noValidation bool) (*Group, error) {
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

	return &Group{
		GroupFormData: f,
		tools: tools{
			log: c.Logger(),
			db:  db,
		},
	}, nil
}

func (t *Group) Insert(ctx context.Context, timeoutSecs int) error {
	t.log.Debugf("group data to insert: %v", t.GroupFormData)

	if errInsert := t.db.Insert(&t.GroupFormData); errInsert != nil {
		return errInsert
	}
	return nil
}
