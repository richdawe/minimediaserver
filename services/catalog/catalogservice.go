package catalog

import (
	"io"

	"github.com/richdawe/minimediaserver/services/storage"
)

type CatalogService interface {
	AddStorage(ss storage.StorageService) error // Add a storage service, its tracks and its playlists to the catalog

	GetTracks() ([]Track, []Playlist)         // Return all the tracks an playlists in the catalog
	GetTrack(id string) (Track, error)        // Get info for a track, by track ID
	ReadTrack(track Track) (io.Reader, error) // Read the track data, using data returned by GetTrack()

	GetPlaylist(id string) (Playlist, error) // Get info for a playlist, by playlist ID
}
