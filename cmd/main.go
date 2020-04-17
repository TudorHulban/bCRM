package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/TudorHulban/bCRM/pkg/constants"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

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
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Routes
	// private routes
	r := e.Group("/r")

	// Start server
	go func() {
		if err := e.Start(constants.ListeningSocket); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	handleInterrupt(e, constants.ShutdownGraceSeconds)
}
