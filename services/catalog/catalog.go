package catalog

import (
	"errors"
	"io"

	"github.com/richdawe/minimediaserver/services/storage"
)

type CatalogService struct {
	storageByID map[string]storage.StorageService // Indexed by storage ID

	// TODO: is this even needed? vvv
	tracksByStorageServiceID map[string][]storage.Track // Indexed by storage ID
	tracksByID               map[string]Track           // Indexed by track ID
	allTracks                []Track
}

func (cs *CatalogService) AddStorage(ss storage.StorageService) error {
	ssid := ss.GetID()
	cs.storageByID[ssid] = ss

	storageTracks := ss.FindTracks()
	cs.tracksByStorageServiceID[ssid] = storageTracks
	for _, storageTrack := range storageTracks {
		track := Track{
			ID:               storageTrack.ID,
			StorageServiceID: ssid,
			Name:             storageTrack.Name,
			MIMEType:         storageTrack.MIMEType,
		}
		cs.tracksByID[track.ID] = track
		cs.allTracks = append(cs.allTracks, track)
	}
	return nil
}

func (cs *CatalogService) GetTracks() []Track {
	return cs.allTracks
}

func (cs *CatalogService) GetTrack(id string) (Track, error) {
	track, ok := cs.tracksByID[id]
	if !ok {
		return Track{}, errors.New("unable to find track by ID")
	}
	return track, nil
}

func (cs *CatalogService) ReadTrack(track Track) (io.Reader, error) {
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

func New() (*CatalogService, error) {
	return &CatalogService{
		storageByID:              make(map[string]storage.StorageService, 0),
		tracksByStorageServiceID: make(map[string][]storage.Track),
		tracksByID:               make(map[string]Track),
	}, nil
}
