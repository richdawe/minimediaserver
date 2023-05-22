package main

import (
	"fmt"
	"os"

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
	e.Logger.Fatal(e.Start(config.Addr))

	fmt.Println("DONE")
}
