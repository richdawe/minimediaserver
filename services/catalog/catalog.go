package catalog

import (
	"github.com/richdawe/minimediaserver/services/storage"
)

type CatalogService struct {
	storages        []storage.StorageService
	tracksByStorage map[string][]storage.Track // Indexed by storage ID
	tracksByID      map[string]Track           // Indexed by track ID
	allTracks       []Track
}

func (cs *CatalogService) AddStorage(s storage.StorageService) error {
	cs.storages = append(cs.storages, s)
	sid := s.GetID()
	storageTracks := s.FindTracks()
	cs.tracksByStorage[sid] = storageTracks
	for _, storageTrack := range storageTracks {
		track := Track{
			StorageService: s,
			ID:             storageTrack.ID,
			Name:           storageTrack.Name,
		}
		cs.tracksByID[track.ID] = track
		cs.allTracks = append(cs.allTracks, track)
	}
	return nil
}

func (cs *CatalogService) GetTracks() []Track {
	return cs.allTracks
}

func New() (*CatalogService, error) {
	return &CatalogService{
		storages:        make([]storage.StorageService, 0),
		tracksByStorage: make(map[string][]storage.Track),
		tracksByID:      make(map[string]Track),
	}, nil
}
