package main

import (
	"fmt"
	"io"
	"os"

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
	fmt.Println("Hai Rich")

	catalogService, err := catalog.New()
	handleErr(err)

	nullStorage, err := storage.NewNullStorage()
	handleErr(err)
	err = catalogService.AddStorage(nullStorage)
	handleErr(err)
	diskStorage, err := storage.NewDiskStorage("Music/cds")
	handleErr(err)
	err = catalogService.AddStorage(diskStorage)
	handleErr(err)

	tracks := catalogService.GetTracks()
	for _, track := range tracks {
		fmt.Println(track.Name, track.ID)
		r, err := track.StorageService.ReadTrack(track.ID)
		if err != nil {
			fmt.Printf("error reading track %s: %s\n", track.ID, err)
			continue
		}

		_, err = io.ReadAll(r)
		if err != nil {
			fmt.Printf("error reading data for track %s: %s\n", track.ID, err)
			continue
		} else {
			fmt.Printf("read track %s\n", track.ID)
		}
	}

	e, err := setupEndpoints(catalogService)
	handleErr(err)
	e.Logger.Fatal(e.Start(":1323"))

	fmt.Println("DONE")
}
