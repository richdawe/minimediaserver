package storage

import "io"

// TODO: what's the proper idiomatic place for an interface
// shared by multiple implementations? in the service pkg, or with where it's used/called?
type StorageService interface {
	GetID() string
	FindTracks() []Track
	ReadTrack(ID string) (io.Reader, error) // may need better name - GetTrack?
}
