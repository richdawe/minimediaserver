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
		tracks := s.FindTracks()
		assert.NotNil(t, tracks)
		assert.Len(t, tracks, 1)
		id := tracks[0].ID
		assert.NotEqual(t, id, "example.ogg")
		assert.Equal(t, tracks[0], Track{
			Name:     "Example",
			ID:       id,
			Location: "/null/example.ogg",
			MIMEType: "audio/ogg",
		})
	})
	t.Run("ReadTrack", func(t *testing.T) {
		tracks := s.FindTracks()
		require.NotNil(t, tracks)
		r, err := s.ReadTrack(tracks[0].ID)
		require.NoError(t, err)

		data, err := io.ReadAll(r)
		require.NoError(t, err)
		dataLen := len(data)
		require.Greater(t, dataLen, 0, "data returned")
	})
}
