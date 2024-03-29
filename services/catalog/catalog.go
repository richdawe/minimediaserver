package catalog

import (
	"errors"
	"io"
	"sort"

	"github.com/richdawe/minimediaserver/services/storage"
)

type BasicCatalog struct {
	storageByID map[string]storage.StorageService // Indexed by storage ID

	// TODO: is this even needed? vvv
	tracksByStorageServiceID    map[string][]storage.Track    // Indexed by storage ID
	playlistsByStorageServiceID map[string][]storage.Playlist // Indexed by storage ID

	tracksByID    map[string]Track    // Indexed by track ID
	playlistsByID map[string]Playlist // Indexed by playlist ID
	allTracks     []Track
	allPlaylists  []Playlist
}

// Assumptions:
// * No ID collisions of tracks from different storage services
func (cs *BasicCatalog) AddStorage(ss storage.StorageService) error {
	ssid := ss.GetID()
	cs.storageByID[ssid] = ss

	storageTracks, storagePlaylists, err := ss.FindTracks()
	if err != nil {
		return err
	}

	cs.tracksByStorageServiceID[ssid] = storageTracks
	for _, storageTrack := range storageTracks {
		track := Track{
			ID:               storageTrack.ID,
			StorageServiceID: ssid,
			Name:             storageTrack.Name,
			MIMEType:         storageTrack.MIMEType,
			DataLen:          storageTrack.DataLen,
		}
		cs.tracksByID[track.ID] = track
		cs.allTracks = append(cs.allTracks, track)
	}

	cs.playlistsByStorageServiceID[ssid] = storagePlaylists
	for _, storagePlaylist := range storagePlaylists {
		playlist := Playlist{
			ID:               storagePlaylist.ID,
			StorageServiceID: ssid,
			Name:             storagePlaylist.Name,
			Tracks:           make([]Track, 0),
		}
		for _, storageTrack := range storagePlaylist.Tracks {
			track := Track{
				ID:               storageTrack.ID,
				StorageServiceID: ssid,
				Name:             storageTrack.Name,
				MIMEType:         storageTrack.MIMEType,
				DataLen:          storageTrack.DataLen,
			}
			playlist.Tracks = append(playlist.Tracks, track)
		}

		cs.playlistsByID[playlist.ID] = playlist
	}
	cs.allPlaylists = sortPlaylists(cs.playlistsByID)

	return nil
}

// Find and sort the playlist IDs based on the name
// of the playlist. Then build list of sorted playlists
// using the sorted list of playlist IDs.
func sortPlaylists(playlistsByID map[string]Playlist) []Playlist {
	playlistIDs := make([]string, 0, 1)
	for _, playlist := range playlistsByID {
		playlistIDs = append(playlistIDs, playlist.ID)
	}
	sort.Slice(playlistIDs, func(i int, j int) bool {
		playlistI := playlistsByID[playlistIDs[i]]
		playlistJ := playlistsByID[playlistIDs[j]]
		return playlistI.Name < playlistJ.Name
	})

	playlists := make([]Playlist, 0, 1)
	for _, playlistID := range playlistIDs {
		playlists = append(playlists, playlistsByID[playlistID])
	}
	return playlists
}

func (cs *BasicCatalog) GetTracks() ([]Track, []Playlist) {
	return cs.allTracks, cs.allPlaylists
}

func (cs *BasicCatalog) GetTrack(id string) (Track, error) {
	track, ok := cs.tracksByID[id]
	if !ok {
		return Track{}, errors.New("unable to find track by ID")
	}
	return track, nil
}

func (cs *BasicCatalog) GetPlaylist(id string) (Playlist, error) {
	playlist, ok := cs.playlistsByID[id]
	if !ok {
		return Playlist{}, errors.New("unable to find playlist by ID")
	}
	return playlist, nil
}

func (cs *BasicCatalog) ReadTrack(track Track) (io.Reader, error) {
	_, err := cs.GetTrack(track.ID)
	if err != nil {
		return nil, err
	}
	ss, ok := cs.storageByID[track.StorageServiceID]
	if !ok {
		return nil, errors.New("unable to find storage service for track")
	}
	return ss.ReadTrack(track.ID)
}

func NewBasicCatalog() (CatalogService, error) {
	return &BasicCatalog{
		storageByID:                 make(map[string]storage.StorageService, 0),
		tracksByStorageServiceID:    make(map[string][]storage.Track),
		playlistsByStorageServiceID: make(map[string][]storage.Playlist),
		tracksByID:                  make(map[string]Track),
		playlistsByID:               make(map[string]Playlist),
	}, nil
}
