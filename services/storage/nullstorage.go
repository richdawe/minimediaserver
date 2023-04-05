package storage

import (
	"bytes"
	"embed"
	"errors"
	"io"

	"github.com/google/uuid"
)

//go:embed example.ogg
var exampleFS embed.FS
var exampleFilename = "example.ogg"

type NullStorage struct {
	ID string

	tracksByID map[string]Track

	Tracks    []Track
	Playlists []Playlist
}

func (ns *NullStorage) GetID() string {
	return ns.ID
}

func (ns *NullStorage) FindTracks() ([]Track, []Playlist, error) {
	if ns.Tracks != nil && ns.Playlists != nil {
		return ns.Tracks, ns.Playlists, nil
	}

	trackLocation := "/null/" + exampleFilename
	trackUUID := locationToUUIDString(trackLocation)
	track := Track{
		Name:     "Example",
		ID:       trackUUID,
		Location: trackLocation,
		MIMEType: "audio/ogg",
	}
	tracks := []Track{track}
	ns.tracksByID[track.ID] = track

	playlist := Playlist{
		Name:     "null-playlist",
		ID:       uuid.New().String(),
		Location: "/null",
		Tracks:   tracks,
	}
	playlists := []Playlist{playlist}

	ns.Tracks = tracks
	ns.Playlists = playlists
	return ns.Tracks, ns.Playlists, nil
}

func (ns *NullStorage) ReadTrack(id string) (io.Reader, error) {
	_, ok := ns.tracksByID[id]
	if !ok {
		// TODO: look at standardizing errors
		return nil, errors.New("track not found")
	}

	data, err := exampleFS.ReadFile(exampleFilename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func NewNullStorage() (*NullStorage, error) {
	return &NullStorage{
		ID:         uuid.New().String(),
		tracksByID: make(map[string]Track),
	}, nil
}
