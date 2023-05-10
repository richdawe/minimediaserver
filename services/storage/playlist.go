package storage

type Playlist struct {
	ID       string // ID unique within this storage service
	Location string // Location within storage service (e.g.: filename, URL)

	Name   string  // Name of the playlist
	Tracks []Track // Tracks in this playlist
}
