package catalog

type Playlist struct {
	ID               string // Unique ID from storage service
	StorageServiceID string // Storage service's ID
	Name             string
	Tracks           []Track
}
