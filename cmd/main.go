package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/TudorHulban/bCRM/pkg/constants"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	_, cancel := context.WithCancel(context.Background()) // creating context for app
	defer cancel()

	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.DEBUG)

	// Start server
	go func() {
		if err := e.Start(constants.ListeningSocket); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	graceSeconds, _ := strconv.Atoi(constants.ListeningSocket)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(graceSeconds)*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
