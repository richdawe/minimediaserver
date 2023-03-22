package storage

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskStorage(t *testing.T) {
	s, err := NewDiskStorage("../../test/services/storage/diskstorage/Music/cds")
	require.NoError(t, err)

	t.Run("GetID", func(t *testing.T) {
		id := s.GetID()
		assert.NotEmpty(t, id)
	})

	t.Run("FindTracks", func(t *testing.T) {
		tracks := s.FindTracks()
		assert.NotNil(t, tracks)
		assert.Equal(t, len(tracks), 3)

		trackIDs := []string{
			tracks[0].ID, tracks[1].ID, tracks[2].ID,
		}

		assert.Equal(t, []Track{
			{
				Name:     "track1-example.ogg",
				ID:       trackIDs[0],
				Location: "../../test/services/storage/diskstorage/Music/cds/Artist/Album1/track1-example.ogg",
				MIMEType: "audio/ogg",
			},
			{
				Name:     "track2-example.flac",
				ID:       trackIDs[1],
				Location: "../../test/services/storage/diskstorage/Music/cds/Artist/Album1/track2-example.flac",
				MIMEType: "audio/flac",
			},
			{
				Name:     "track1-example.ogg",
				ID:       trackIDs[2],
				Location: "../../test/services/storage/diskstorage/Music/cds/Artist/Album2/track1-example.ogg",
				MIMEType: "audio/ogg",
			},
		}, tracks)

		assert.NotEmpty(t, trackIDs[0])
		assert.NotEmpty(t, trackIDs[1])
		assert.NotEmpty(t, trackIDs[2])
		assert.NotEqual(t, trackIDs[0], trackIDs[1])
		assert.NotEqual(t, trackIDs[1], trackIDs[2])
	})

	t.Run("ReadTrack", func(t *testing.T) {
		tracks := s.FindTracks()
		require.NotNil(t, tracks)
		assert.Equal(t, len(tracks), 3)

		dataLens := make([]int, 0)

		for _, track := range tracks {
			r, err := s.ReadTrack(track.ID)
			require.NoError(t, err)

			data, err := io.ReadAll(r)
			require.NoError(t, err)
			dataLens = append(dataLens, len(data))
		}

		// Paths below relative to test/services/storage/diskstorage
		require.Equal(t, []int{
			104793, // Music/cds/Artist/Album1/track1-example.ogg
			980027, // Music/cds/Artist/Album1/track2-example.flac
			104793, // Music/cds/Artist/Album2/track1-example.ogg
		}, dataLens)
	})
}
