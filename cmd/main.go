package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/TudorHulban/bCRM/models"
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var dbConn *pg.DB

func main() {
	dbconnInfo := DBConnInfo{
		Socket: commons.DBSocket,
		User:   commons.DBUser,
		Pass:   commons.DBPass,
		DB:     commons.DBName,
	}
	dbConn := pg.Connect(&pg.Options{
		Addr:     dbconnInfo.Socket,
		User:     dbconnInfo.User,
		Password: dbconnInfo.Pass,
		Database: dbconnInfo.DB,
	})
	defer dbConn.Close()

	// check db conn
	errConn := CheckPgDB(dbConn)
	if errConn != nil {
		log.Print("Could not create DB schema. Exiting ...", errConn)
		os.Exit(1)
	}

	// Create DB schema
	errSchema := NewSchema(dbConn, interface{}(&models.SLAPriority{}), interface{}(&models.SLA{}), interface{}(&models.SLAValue{}), interface{}(&models.TicketType{}), interface{}(&models.TicketStatus{}), interface{}(&models.Resource{}), interface{}(&models.ResourceMove{}), interface{}(&models.Event{}), interface{}(&models.TicketMovement{}), interface{}(&models.Ticket{}), interface{}(&models.Team{}), interface{}(&models.User{}), interface{}(&models.File{}), interface{}(&models.Contact{}))
	if errSchema != nil {
		log.Print("Could not create DB schema. Exiting ...", errSchema)
		os.Exit(1)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Routes
	// Public routes
	e.GET(commons.EndpointLive, Live)
	//e.POST(EndpointLogin, LoginWithPassword)
	e.POST(commons.EndpointNewUser, NewUser)

	// private routes
	//r := e.Group("/r")

	// Start server
	go func() {
		if err := e.Start(commons.ListeningSocket); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	handleInterrupt(e, commons.ShutdownGraceSeconds)
}

func handleInterrupt(s *echo.Echo, graceSeconds int) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint // the read from channel blocks until interrupt is received and sent on channel.

	// we can now shutdown
	log.Print("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(graceSeconds)*time.Second)
	defer cancel()

	if errShutdown := s.Shutdown(ctx); errShutdown != nil {
		log.Printf("Error HTTP server shutdown: %v", errShutdown)
	}
	time.Sleep(time.Duration(graceSeconds))
}
