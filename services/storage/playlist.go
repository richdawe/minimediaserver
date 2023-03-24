package storage

type Playlist struct {
	Name     string  // Name of the playlist
	ID       string  // ID unique within this storage service
	Location string  // Location within storage service (e.g.: filename, URL)
	Tracks   []Track // Tracks in this playlist
}
