package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4/middleware"

	"github.com/richdawe/minimediaserver/services/catalog"
	"github.com/richdawe/minimediaserver/services/storage"
)

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	catalogService, err := catalog.New()
	handleErr(err)

	// TODO: need a config file for configuring storage backends
	nullStorage, err := storage.NewNullStorage()
	handleErr(err)
	err = catalogService.AddStorage(nullStorage)
	handleErr(err)
	diskStorage, err := storage.NewDiskStorage("test/services/storage/diskstorage")
	handleErr(err)
	err = catalogService.AddStorage(diskStorage)
	handleErr(err)

	e, err := setupEndpoints(catalogService)
	handleErr(err)

	// TODO: need a config file for specifying HTTP server options
	e.Use(middleware.Timeout())
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))

	fmt.Println("DONE")
}
