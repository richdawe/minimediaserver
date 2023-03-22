package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
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

// Find the tracks in this storage, and return the tracks
// in a stable order.
func (ds *DiskStorage) FindTracks() []Track {
	if ds.tracksByID == nil {
		ds.populateTracks()
	}

	// Find and sort the track IDs based on the location
	// of the track.
	trackIDs := make([]string, 0, 1)
	for _, track := range ds.tracksByID {
		trackIDs = append(trackIDs, track.ID)
	}
	sort.Slice(trackIDs, func(i int, j int) bool {
		trackI := ds.tracksByID[trackIDs[i]]
		trackJ := ds.tracksByID[trackIDs[j]]
		return trackI.Location < trackJ.Location
	})

	tracks := make([]Track, 0, 1)
	for _, trackID := range trackIDs {
		tracks = append(tracks, ds.tracksByID[trackID])
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

func ignoreMIMEType(mimeType string) bool {
	switch mimeType {
	case "application/binary":
		return true
	}
	return false
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
		location := ds.BasePath + "/" + path // TODO: proper cross platform concat

		// Ignore some unknown MIME types
		mimeType := getMIMEType(d.Name())
		if ignoreMIMEType(mimeType) {
			fmt.Printf("Ignoring file due to MIME type: %s\n", location)
			return nil
		}

		track := Track{
			Name:     d.Name(),
			ID:       uuid.New().String(),
			Location: location,
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
