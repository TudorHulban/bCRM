package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/TudorHulban/bCRM/models"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

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

	// check db conn
	errConn := commons.CheckPgDB(e.Logger)
	if errConn != nil {
		log.Print("Could not create DB schema. Exiting ...", errConn)
		os.Exit(1)
	}

	// check schema was created already
	if exists, errExists := models.TableExists(ctx, &models.UserData{}, "users", 5, e.Logger); errExists != nil {
		e.Logger.Debug("table users does not exist:", errExists, exists)
		// currently does not work. issue has been opened in orm lib
	}

	// create DB schema
	if errSchema := createSchema(commons.DB()); errSchema != nil {
		log.Print("Could not create DB schema. Exiting ...", errSchema)
		os.Exit(1)
	}

	// populate schema
	if errPopul := populateSchema(ctx, e.NewContext(nil, nil), commons.DB()); errPopul != nil {
		log.Print("Could not populate DB schema. Exiting ...", errPopul)
		os.Exit(1)
	}

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
	commons.DB().Close()
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
