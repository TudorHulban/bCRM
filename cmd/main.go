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
var ctx context.Context

func main() {
	ctx = context.Background()

	e := echo.New()
	e.HideBanner = true
	e.DisableHTTP2 = true
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	dbConn = pg.Connect(&pg.Options{
		Addr:     commons.DBSocket,
		User:     commons.DBUser,
		Password: commons.DBPass,
		Database: commons.DBName,
	})

	// check db conn
	errConn := commons.CheckPgDB(e.Logger, dbConn)
	if errConn != nil {
		log.Print("Could not create DB schema. Exiting ...", errConn)
		os.Exit(1)
	}

	// check schema was created already
	exists, errExists := models.TableExists(ctx, dbConn, &models.UserData{}, "users", 5, e.Logger)
	if errExists != nil {
		e.Logger.Debug("table users does not exist:", errExists, exists)
		// currently does not work. issue has been opened in orm lib
	}

	// Create DB schema
	errSchema := createSchema(dbConn)
	if errSchema != nil {
		log.Print("Could not create DB schema. Exiting ...", errSchema)
		os.Exit(1)
	}

	// populate schema
	populateSchema(ctx, e.NewContext(nil, nil), dbConn)

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
	e.Logger.Info("closing DB")
	dbConn.Close() // to switch to defer maybe
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
