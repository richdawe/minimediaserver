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
}

func (ns *NullStorage) GetID() string {
	return ns.ID
}

func (ns *NullStorage) FindTracks() []Track {
	tracks := make([]Track, 0, 1)
	track := Track{
		Name:     "Example",
		ID:       uuid.New().String(),
		Location: "/null/" + exampleFilename,
		MIMEType: "audio/ogg",
	}
	tracks = append(tracks, track)
	return tracks
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
