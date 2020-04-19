package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/TudorHulban/bCRM/constants"
	"github.com/TudorHulban/bCRM/interfaces"
	"github.com/TudorHulban/bCRM/pgstore"
	"github.com/TudorHulban/bCRM/pkg/httphandlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	dbconnInfo := interfaces.DBConnInfo{
		Socket: constants.DBSocket,
		User:   constants.DBUser,
		Pass:   constants.DBPass,
		DB:     constants.DBName,
	}
	var store pgstore.PgStore
	db, err := store.Open(dbconnInfo)
	if err != nil {
		log.Print("Could not connect to DB. Exiting ...")
		os.Exit(1)
	}
	// set schema

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Routes
	// public routes
	e.GET(constants.EndpointLive, httphandlers.Live)
	e.POST(constants.EndpointLogin, httphandlers.LoginWithPassword)
	e.POST(constants.EndpointNewUser, httphandlers.NewUser)

	// private routes
	//r := e.Group("/r")

	// Start server
	go func() {
		if err := e.Start(constants.ListeningSocket); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	handleInterrupt(e, constants.ShutdownGraceSeconds)
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
