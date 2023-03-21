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
		assert.Equal(t, tracks[0], Track{
			Name: "Example",
			ID:   "example.ogg",
		})
	})
	t.Run("ReadTrack", func(t *testing.T) {
		r, err := s.ReadTrack("example.ogg")
		require.NoError(t, err)

		data, err := io.ReadAll(r)
		require.NoError(t, err)
		dataLen := len(data)
		require.Greater(t, dataLen, 0, "data returned")
	})
}