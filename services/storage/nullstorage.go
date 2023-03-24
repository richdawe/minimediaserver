package storage

import (
	"bytes"
	"embed"
	"io"

	"github.com/google/uuid"
)

//go:embed example.ogg
var exampleFS embed.FS
var exampleFilename = "example.ogg"

type NullStorage struct {
	ID string

	Tracks    []Track
	Playlists []Playlist
}

func (ns *NullStorage) GetID() string {
	return ns.ID
}

func (ns *NullStorage) FindTracks() ([]Track, []Playlist) {
	if ns.Tracks != nil && ns.Playlists != nil {
		return ns.Tracks, ns.Playlists
	}

	tracks := make([]Track, 0, 1)
	track := Track{
		Name:     "Example",
		ID:       uuid.New().String(),
		Location: "/null/" + exampleFilename,
		MIMEType: "audio/ogg",
	}
	tracks = append(tracks, track)

	playlist := Playlist{
		Name:     "null-playlist",
		ID:       uuid.New().String(),
		Location: "/null",
		Tracks:   tracks,
	}
	playlists := []Playlist{playlist}

	ns.Tracks = tracks
	ns.Playlists = playlists
	return ns.Tracks, ns.Playlists
}

func (ns *NullStorage) ReadTrack(id string) (io.Reader, error) {
	data, err := exampleFS.ReadFile(exampleFilename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func NewNullStorage() (*NullStorage, error) {
	return &NullStorage{
		ID: uuid.New().String(),
	}, nil
}
