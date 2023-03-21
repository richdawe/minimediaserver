package storage

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/google/uuid"
)

type DiskStorage struct {
	ID   string
	Path string
}

func (ds *DiskStorage) GetID() string {
	return ds.ID
}

func (ds *DiskStorage) FindTracks() []Track {
	tracks := make([]Track, 0, 1)

	fileSystem := os.DirFS(ds.Path)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if d.IsDir() {
			return nil
		}
		// TODO: ignore some file extensions

		fmt.Println(d.Name(), path)
		track := Track{
			Name: d.Name(),
			ID:   ds.Path + "/" + path, // TODO: proper concat
		}
		tracks = append(tracks, track)
		return nil
	})

	return tracks
}

func (ds *DiskStorage) ReadTrack(ID string) (io.Reader, error) {
	// TODO: Some way to not expose paths as IDs? E.g.: generate a UUID for each thing in here,
	// then map UUID -> name, path. Having paths in public API is going to leave this open
	// to (potentially) path traversal attacks.

	// TODO: ID needs some kind of validation here!

	// TODO: figure out some way to return a reader that reads in the file in chunks
	data, err := os.ReadFile(ID)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func NewDiskStorage(Path string) (*DiskStorage, error) {
	return &DiskStorage{
		ID:   uuid.New().String(),
		Path: Path,
	}, nil
}
