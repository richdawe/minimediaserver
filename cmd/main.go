package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4/middleware"
)

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	setLoadConfigOptions()
	config, err := loadConfig()
	handleErr(err)

	catalogService, err := buildCatalog(config)
	handleErr(err)

	e, err := setupEndpoints(config, catalogService)
	handleErr(err)

	// TODO: need a config file for specifying HTTP server options
	e.Use(middleware.Timeout())
	e.Use(middleware.Logger())

	// Start server
	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	// https://medium.com/@mokiat/proper-http-shutdown-in-go-bd3bfaade0f2
	go func() {
		if err := e.Start(config.Addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(fmt.Sprintf("shutting down the server: %s", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
		e.Close()
	}

	fmt.Println("DONE")
}
