package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

type DiskStorage struct {
	ID       string
	BasePath string

	tracksByID map[string]Track
}

func (ds *DiskStorage) GetID() string {
	return ds.ID
}

func (ds *DiskStorage) FindTracks() []Track {
	if ds.tracksByID == nil {
		ds.populateTracks()
	}
	tracks := make([]Track, 0, 1)
	for _, track := range ds.tracksByID {
		tracks = append(tracks, track)
	}
	return tracks
}

func getMIMEType(filename string) string {
	var mimeType string

	filename = strings.ToLower(filename)
	switch {
	case strings.HasSuffix(filename, ".mp3"):
		mimeType = "audio/mp3"
	case strings.HasSuffix(filename, ".ogg"):
		mimeType = "audio/ogg"
	case strings.HasSuffix(filename, ".flac"):
		mimeType = "audio/flac"
	}

	if mimeType == "" {
		mimeType = "application/binary"
	}
	return mimeType
}

func (ds *DiskStorage) populateTracks() {
	tracks := make(map[string]Track, 0)

	fileSystem := os.DirFS(ds.BasePath)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if d.IsDir() {
			return nil
		}
		// TODO: ignore some file extensions and MIME types
		mimeType := getMIMEType(d.Name())

		fmt.Println(d.Name(), path)
		track := Track{
			Name:     d.Name(),
			ID:       uuid.New().String(),
			Location: ds.BasePath + "/" + path, // TODO: proper cross platform concat
			MIMEType: mimeType,
		}
		tracks[track.ID] = track
		return nil
	})

	ds.tracksByID = tracks
}

func (ds *DiskStorage) ReadTrack(id string) (io.Reader, error) {
	track, ok := ds.tracksByID[id]
	if !ok {
		// TODO: look at standardizing errors
		return nil, errors.New("track not found")
	}

	// TODO: figure out some way to return a reader that reads in the file in chunks
	data, err := os.ReadFile(track.Location)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func NewDiskStorage(Path string) (*DiskStorage, error) {
	return &DiskStorage{
		ID:       uuid.New().String(),
		BasePath: Path,
	}, nil
}
