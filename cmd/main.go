package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/TudorHulban/bCRM/constants"
	"github.com/TudorHulban/bCRM/pkg/httphandlers"
	"github.com/TudorHulban/bCRM/variables"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// init
	initStore(variables.GStore)

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
	e.GET(constants.EndpointOK, httphandlers.HandlerOK)
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
