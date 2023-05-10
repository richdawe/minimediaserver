package storage

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullStorage(t *testing.T) {
	s, err := NewNullStorage()
	require.NoError(t, err)

	t.Run("GetID", func(t *testing.T) {
		id := s.GetID()
		assert.NotEmpty(t, id)
	})

	t.Run("FindTracks", func(t *testing.T) {
		tracks, playlists, err := s.FindTracks()
		require.NoError(t, err)
		assert.NotNil(t, tracks)
		assert.Len(t, tracks, 1)

		assert.NotNil(t, playlists)
		assert.Len(t, playlists, 1)

		id := tracks[0].ID
		assert.NotEqual(t, id, "example.ogg")
		assert.Equal(t, tracks[0], Track{
			// TITLE comment checked and set using:
			// vorbiscomment services/storage/example.ogg
			// vorbiscomment -a -t title=ExAmPlE services/storage/example.ogg
			Name:     "ExAmPlE",
			ID:       id,
			Location: "/null/example.ogg",
			MIMEType: "audio/ogg",
			DataLen:  105269,
		})

		playlistID := playlists[0].ID
		assert.NotEqual(t, playlistID, "example.ogg")
		assert.NotEqual(t, playlistID, id)
		assert.Equal(t, playlists[0], Playlist{
			Name:     "null-playlist",
			ID:       playlistID,
			Location: "/null",
			Tracks:   []Track{tracks[0]},
		})
	})

	t.Run("StableIDs", func(t *testing.T) {
		// Verify that the track ID and playlist ID are stable across
		// calls to FindTracks.
		tracks, playlists, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		require.Equal(t, len(tracks), 1)
		require.NotNil(t, playlists)
		require.Equal(t, len(playlists), 1)

		tracks2, playlists2, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		require.Equal(t, len(tracks), 1)
		require.NotNil(t, playlists)
		require.Equal(t, len(playlists), 1)

		assert.Equal(t, tracks, tracks2)
		assert.Equal(t, playlists, playlists2)
	})

	t.Run("ReadTrack", func(t *testing.T) {
		tracks, _, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		r, err := s.ReadTrack(tracks[0].ID)
		require.NoError(t, err)

		data, err := io.ReadAll(r)
		require.NoError(t, err)
		dataLen := len(data)
		require.Greater(t, dataLen, 0, "data returned")
	})

	t.Run("ReadTrackNotFound", func(t *testing.T) {
		_, err := s.ReadTrack("nope")
		require.Error(t, err)
	})
}
