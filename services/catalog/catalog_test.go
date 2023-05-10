package catalog

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		assert.Equal(t, Track{
			Name:             "ExAmPlE",
			ID:               id,
			StorageServiceID: nullStorage.GetID(),
			MIMEType:         "audio/ogg",
			DataLen:          105269,
		}, tracks[0])

		playlistID := playlists[0].ID
		assert.NotEqual(t, playlistID, "example.ogg")
		assert.NotEqual(t, playlistID, id)
		assert.Equal(t, Playlist{
			Name:             "null-playlist",
			ID:               playlistID,
			StorageServiceID: nullStorage.GetID(),
			Tracks:           []Track{tracks[0]},
		}, playlists[0])
	})

	t.Run("GetTrack", func(t *testing.T) {
		tracks, _ := catalogService.GetTracks()
		require.NotNil(t, tracks)
		require.Len(t, tracks, 1)

		track, err := catalogService.GetTrack(tracks[0].ID)
		require.NoError(t, err)
		assert.Equal(t, tracks[0], track)

		_, err = catalogService.GetTrack("nope")
		assert.Error(t, err)
	})

	t.Run("GetPlaylist", func(t *testing.T) {
		_, playlists := catalogService.GetTracks()
		require.NotNil(t, playlists)
		require.NoError(t, err)

		playlist, err := catalogService.GetPlaylist(playlists[0].ID)
		require.NoError(t, err)
		assert.Equal(t, playlists[0], playlist)

		_, err = catalogService.GetPlaylist("nope")
		assert.Error(t, err)
	})

	t.Run("ReadTrack", func(t *testing.T) {
		tracks, _ := catalogService.GetTracks()
		require.NotNil(t, tracks)
		require.Len(t, tracks, 1)

		r, err := catalogService.ReadTrack(tracks[0])
		require.NoError(t, err)

		data, err := io.ReadAll(r)
		require.NoError(t, err)
		dataLen := len(data)
		require.Greater(t, dataLen, 0, "data returned")

		// Invalid track ID
		track := Track{ID: "nope", MIMEType: "audio/nope"}
		_, err = catalogService.ReadTrack(track)
		require.Error(t, err)

		// Invalid storage service ID
		track = tracks[0]
		track.StorageServiceID = "nope"
		_, err = catalogService.ReadTrack(track)
		require.Error(t, err)
	})
}
