package storage

import "io"

type StorageService interface {
	GetID() string

	FindTracks() ([]Track, []Playlist, error)
	ReadTrack(id string) (io.Reader, error) // may need better name - GetTrack?
}
