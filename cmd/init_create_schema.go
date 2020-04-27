package main

import (
	"github.com/TudorHulban/bCRM/models"
	"github.com/go-pg/pg/v9"
)

func createSchema(db *pg.DB) error {
	tables := []interface{}{}

	tables = append(tables, interface{}(&models.TeamFormData{}))
	tables = append(tables, interface{}(&models.GroupFormData{}))
	tables = append(tables, interface{}(&models.UserData{}))

	return newSchema(db, tables)
}

// interface{}(&models.SLAPriority{}), interface{}(&models.SLA{}), interface{}(&models.SLAValue{}), interface{}(&models.TicketType{}),
// interface{}(&models.TicketStatus{}), interface{}(&models.Resource{}), interface{}(&models.ResourceMove{}), interface{}(&models.Event{}),
// interface{}(&models.TicketMovement{}), interface{}(&models.Ticket{}), interface{}(&models.File{}), interface{}(&models.Contact{}))
//
