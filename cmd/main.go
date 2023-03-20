package main

import (
	"fmt"
	"io"
	"os"

	"github.com/richdawe/minimediaserver/services/storage"
)

type StorageService interface {
	FindTracks() []storage.Track
	ReadTrack(ID string) (io.Reader, error) // may need better name - GetTrack?
}

func main() {
	fmt.Println("Hai Rich")

	var storageService StorageService
	//storageService, err := storage.NewNullStorage()
	storageService, err := storage.NewDiskStorage("Music/cds")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	tracks := storageService.FindTracks()
	for _, track := range tracks {
		fmt.Println(track.Name, track.ID)
		r, err := storageService.ReadTrack(track.ID)
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

	fmt.Println("DONE")
}
