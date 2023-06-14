package storage

import (
	"io"
	"io/fs"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskStorage(t *testing.T) {
	s, err := NewDiskStorage("../../testdata/services/storage/diskstorage/Music/cds")
	require.NoError(t, err)

	t.Run("GetID", func(t *testing.T) {
		id := s.GetID()
		assert.NotEmpty(t, id)
	})

	t.Run("FindTracks", func(t *testing.T) {
		tracks, playlists, err := s.FindTracks()
		require.NoError(t, err)
		assert.NotNil(t, tracks)
		assert.Equal(t, len(tracks), 4)

		assert.NotNil(t, playlists)
		assert.Equal(t, len(playlists), 2)

		// Check tracks
		trackIDs := []string{
			tracks[0].ID, tracks[1].ID, tracks[2].ID, tracks[3].ID,
		}

		assert.Equal(t, []Track{
			{
				Name:     "ALBUM1_TRACK1_EXAMPLE",
				Tags:     Tags{Title: "ALBUM1_TRACK1_EXAMPLE", Album: "album1", Artist: "the-artist", Genre: "Example", TrackNumber: 1},
				ID:       trackIDs[0],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album1/track1-example.ogg",
				MIMEType: "audio/ogg",
				DataLen:  105354,
			},
			{
				Name:     "ALBUM1_TRACK2_EXAMPLE",
				Tags:     Tags{Title: "ALBUM1_TRACK2_EXAMPLE", Album: "album1", Artist: "the-artist", Genre: "ExampleMulti-value", TrackNumber: 2},
				ID:       trackIDs[1],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album1/track2-example.flac",
				MIMEType: "audio/flac",
				DataLen:  980027,
			},
			{
				Name:     "ALBUM2_TRACK1_EXAMPLE",
				Tags:     Tags{Title: "ALBUM2_TRACK1_EXAMPLE", Album: "album2", Artist: "another-artist"},
				ID:       trackIDs[2],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album2/track1-example.ogg",
				MIMEType: "audio/ogg",
				DataLen:  105324,
			},
			// TODO: fix tags
			{
				Name:     "ALBUM1_TRACK2_EXAMPLE\x00",
				Tags:     Tags{Title: "ALBUM1_TRACK2_EXAMPLE\x00", Album: "album1\x00", Artist: "the-artist\x00", Genre: "Example;Multi-value\x00"},
				ID:       trackIDs[3],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album2/track2-example.mp3",
				MIMEType: "audio/mp3",
				DataLen:  161632,
			},
		}, tracks)

		assert.NotEmpty(t, trackIDs[0])
		assert.NotEmpty(t, trackIDs[1])
		assert.NotEmpty(t, trackIDs[2])
		assert.NotEmpty(t, trackIDs[3])
		assert.NotEqual(t, trackIDs[0], trackIDs[1])
		assert.NotEqual(t, trackIDs[1], trackIDs[2])
		assert.NotEqual(t, trackIDs[2], trackIDs[3])

		// Check playlists
		playlistIDs := []string{
			playlists[0].ID, playlists[1].ID,
		}

		assert.Equal(t, []Playlist{
			{
				Name:     "Artist :: Album1",
				ID:       playlistIDs[0],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album1",
				Tracks:   []Track{tracks[0], tracks[1]},
			},
			{
				Name:     "Artist :: Album2",
				ID:       playlistIDs[1],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album2",
				Tracks:   []Track{tracks[2], tracks[3]},
			},
		}, playlists)
	})

	t.Run("StableIDs", func(t *testing.T) {
		// Verify that the track ID and playlist ID are stable across
		// calls to FindTracks.
		tracks, playlists, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		require.Equal(t, len(tracks), 4)
		require.NotNil(t, playlists)
		require.Equal(t, len(playlists), 2)

		tracks2, playlists2, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		require.Equal(t, len(tracks), 4)
		require.NotNil(t, playlists)
		require.Equal(t, len(playlists), 2)

		assert.Equal(t, tracks, tracks2)
		assert.Equal(t, playlists, playlists2)
	})

	t.Run("ReadTrack", func(t *testing.T) {
		tracks, _, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		assert.Equal(t, len(tracks), 4)

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
			105354, // Music/cds/Artist/Album1/track1-example.ogg
			980027, // Music/cds/Artist/Album1/track2-example.flac
			105324, // Music/cds/Artist/Album2/track1-example.ogg
			161632, // Music/cds/Artist/Album2/track2-example.mp3
		}, dataLens)
	})

	t.Run("ReadTrackNotFound", func(t *testing.T) {
		_, err := s.ReadTrack("nope")
		require.Error(t, err)
	})
}

func TestDiskStorageFailures(t *testing.T) {
	t.Run("BadPath", func(t *testing.T) {
		s, err := NewDiskStorage("./__DOES_NOT_EXIST__")
		require.Error(t, err)
		var pErr *fs.PathError
		require.ErrorAs(t, err, &pErr)
		assert.Equal(t, pErr.Err, syscall.ENOENT)
		require.Nil(t, s)
	})
}
