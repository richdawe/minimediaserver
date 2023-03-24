package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/richdawe/minimediaserver/services/storage"
)

func TestCatalogService(t *testing.T) {
	catalogService, err := New()
	assert.Nil(t, err)
	nullStorage, err := storage.NewNullStorage()
	assert.Nil(t, err)
	err = catalogService.AddStorage(nullStorage)
	assert.Nil(t, err)

	t.Run("Loaded NullStorage data into catalog", func(t *testing.T) {
		tracks, playlists := catalogService.GetTracks()
		assert.NotNil(t, tracks)
		assert.Len(t, tracks, 1)

		assert.NotNil(t, playlists)
		assert.Len(t, playlists, 1)

		id := tracks[0].ID
		assert.NotEqual(t, id, "example.ogg")
		assert.Equal(t, tracks[0], Track{
			Name:             "Example",
			ID:               id,
			StorageServiceID: nullStorage.GetID(),
			MIMEType:         "audio/ogg",
		})

		playlistID := playlists[0].ID
		assert.NotEqual(t, playlistID, "example.ogg")
		assert.NotEqual(t, playlistID, id)
		assert.Equal(t, playlists[0], Playlist{
			Name:             "null-playlist",
			ID:               playlistID,
			StorageServiceID: nullStorage.GetID(),
			Tracks:           []Track{tracks[0]},
		})
	})
}
