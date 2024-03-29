package storage

import (
	"io"
	"io/fs"
	"path/filepath"
	"regexp/syntax"
	"syscall"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mp3Regexps = []string{
	// album-artist, album (trackno), artist, title - separate by dashes
	// commented out version doesn't work when album/title includes dashes, so just use whitespace around dash as delimiter
	//"(?P<albumartist>[^-]+) - (?P<album>[^-]+) \\((?P<trackno>\\d+)\\) - (?P<artist>[^-]+) - (?P<title>.+)",
	"(?P<albumartist>.+) - (?P<album>.+) \\((?P<trackno>\\d+)\\) - (?P<artist>.+) - (?P<title>.+)",

	// album, artist (trackno), title - separate by dashes
	// commented out version doesn't work when album/title includes dashes, so just use whitespace around dash as delimiter
	//"(?P<artist>[^-]+) - (?P<album>[^-]+) \\((?P<trackno>\\d+)\\) - (?P<title>.+)",
	"(?P<artist>.+) - (?P<album>.+) \\((?P<trackno>\\d+)\\) - (?P<title>.+)",
}

func TestDiskStorage(t *testing.T) {
	s, err := NewDiskStorage("../../testdata/services/storage/diskstorage/Music/cds", []string{})
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
		assert.Equal(t, len(playlists), 3)

		// Check tracks
		trackIDs := []string{
			tracks[0].ID, tracks[1].ID, tracks[2].ID, tracks[3].ID,
		}

		assert.Equal(t, []Track{
			{
				Name:        "the-artist :: ALBUM1_TRACK1_EXAMPLE",
				Title:       "ALBUM1_TRACK1_EXAMPLE",
				Artist:      "the-artist",
				Album:       "album1",
				AlbumArtist: "Artist",

				PlaylistLocation: "tags:../../testdata/services/storage/diskstorage/Music/cds/Artist/album1",

				Tags: Tags{Title: "ALBUM1_TRACK1_EXAMPLE", Album: "album1", Artist: "the-artist", Genre: "Example", TrackNumber: 1},

				ID:       trackIDs[0],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album1/track1-example.ogg",
				MIMEType: "audio/ogg",
				DataLen:  105354,
			},
			{
				Name:        "the-artist :: ALBUM1_TRACK2_EXAMPLE",
				Title:       "ALBUM1_TRACK2_EXAMPLE",
				Artist:      "the-artist",
				Album:       "album1",
				AlbumArtist: "Artist",

				PlaylistLocation: "tags:../../testdata/services/storage/diskstorage/Music/cds/Artist/album1",

				Tags: Tags{Title: "ALBUM1_TRACK2_EXAMPLE", Album: "album1", Artist: "the-artist", Genre: "ExampleMulti-value", TrackNumber: 2},

				ID:       trackIDs[1],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album1/track2-example.flac",
				MIMEType: "audio/flac",
				DataLen:  980027,
			},
			{
				Name:        "another-artist :: ALBUM2_TRACK1_EXAMPLE",
				Title:       "ALBUM2_TRACK1_EXAMPLE",
				Artist:      "another-artist",
				Album:       "album2",
				AlbumArtist: "Artist",

				PlaylistLocation: "tags:../../testdata/services/storage/diskstorage/Music/cds/Artist/album2",

				Tags: Tags{Title: "ALBUM2_TRACK1_EXAMPLE", Album: "album2", Artist: "another-artist"},

				ID:       trackIDs[2],
				Location: "../../testdata/services/storage/diskstorage/Music/cds/Artist/Album2/track1-example.ogg",
				MIMEType: "audio/ogg",
				DataLen:  105324,
			},
			// TODO: fix tags to not contain nuls - may need a changeset from a PR on id3-go
			{
				Name:        "ALBUM1_TRACK2_EXAMPLE\x00",
				Title:       "ALBUM1_TRACK2_EXAMPLE\x00",
				Artist:      "the-artist\x00",
				Album:       "album1\x00",
				AlbumArtist: "the-artist\x00",

				PlaylistLocation: "tags:../../testdata/services/storage/diskstorage/Music/cds/the-artist\x00/album1\x00",

				Tags: Tags{Title: "ALBUM1_TRACK2_EXAMPLE\x00", Album: "album1\x00", AlbumArtist: "the-artist\x00", Artist: "the-artist\x00", Genre: "Example;Multi-value\x00"},

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
			playlists[0].ID, playlists[1].ID, playlists[2].ID,
		}

		assert.Equal(t, []Playlist{
			{
				Name:     "Artist :: album1",
				ID:       playlistIDs[0],
				Location: "tags:../../testdata/services/storage/diskstorage/Music/cds/Artist/album1",
				Tracks:   []Track{tracks[0], tracks[1]},
			},
			{
				Name:     "Artist :: album2",
				ID:       playlistIDs[1],
				Location: "tags:../../testdata/services/storage/diskstorage/Music/cds/Artist/album2",
				Tracks:   []Track{tracks[2]},
			},
			{
				Name:     "the-artist\x00 :: album1\x00",
				ID:       playlistIDs[2],
				Location: "tags:../../testdata/services/storage/diskstorage/Music/cds/the-artist\x00/album1\x00",
				Tracks:   []Track{tracks[3]},
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
		require.Equal(t, len(playlists), 3)

		tracks2, playlists2, err := s.FindTracks()
		require.NoError(t, err)
		require.NotNil(t, tracks)
		require.Equal(t, len(tracks), 4)
		require.NotNil(t, playlists)
		require.Equal(t, len(playlists), 3)

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
		s, err := NewDiskStorage("./__DOES_NOT_EXIST__", []string{})
		require.Error(t, err)
		var pErr *fs.PathError
		require.ErrorAs(t, err, &pErr)
		assert.Equal(t, pErr.Err, syscall.ENOENT)
		require.Nil(t, s)
	})

	t.Run("BadRegexps", func(t *testing.T) {
		s, err := NewDiskStorage("../../testdata/services/storage/diskstorage/Music/cds", []string{
			"(?P<incomplete", "(?Pperllooking[^-]+)",
		})
		require.Error(t, err)
		var pErr *syntax.Error
		require.ErrorAs(t, err, &pErr)
		require.Nil(t, s)
	})
}

func TestIsTrackByAlbumArtist(t *testing.T) {
	testCases := []struct {
		TrackArtist string
		AlbumArtist string
		Expected    bool
	}{
		// Basic case-insensitivity
		{"Fred", "Fred", true},
		{"Fred", "fred", true},
		{"fred", "FRED", true},
		{"Fred Bloggs", "Fred Bloggs", true},

		// Underscore matching (converted to space)
		{"Fred_Bloggs", "Fred Bloggs", true},
		{"Fred Bloggs", "Fred_Bloggs", true},
		{"Fred_Bloggs", "Fred_Bloggs", true},
		{"Fred Bloggs", "Fred___Bloggs", false},

		// Dash matching (converted to space)
		{"Fred-Bloggs", "Fred Bloggs", true},
		{"Fred Bloggs", "Fred-Bloggs", true},
		{"Fred-Bloggs", "Fred-Bloggs", true},
		{"Fred Bloggs", "Fred---Bloggs", false},
	}

	for i, testCase := range testCases {
		res := isTrackByAlbumArtist(testCase.TrackArtist, testCase.AlbumArtist)
		assert.Equal(t, testCase.Expected, res, "[i=%d] %s %s", i, testCase.TrackArtist, testCase.AlbumArtist)
	}
}

func TestMatchLocation(t *testing.T) {
	testCases := []struct {
		Location      string
		ExpectedTrack Track
		ExpectMatch   bool
	}{
		{
			Location: "adam f - colours (08) - circles.mp3",
			ExpectedTrack: Track{
				Artist: "adam f", Album: "colours", TrackNumber: 8, Title: "circles",
			},
			ExpectMatch: true,
		},
		{
			Location: "a - hi-fi serious (12) - hi-fi serious.mp3",
			ExpectedTrack: Track{
				Artist: "a", Album: "hi-fi serious", TrackNumber: 12, Title: "hi-fi serious",
			},
			ExpectMatch: true,
		},
		{
			Location: "botchit & scarper - botchit breaks (disc 01) (07) - full moon scientist - we demand a shrubbery (scissorkicks mix).mp3",
			ExpectedTrack: Track{
				AlbumArtist: "botchit & scarper",
				Album:       "botchit breaks (disc 01)",
				TrackNumber: 7,
				Artist:      "full moon scientist",
				Title:       "we demand a shrubbery (scissorkicks mix)",
			},
			ExpectMatch: true,
		},
		{
			Location: "arcade fire - funeral (07) - arcade fire - wake up.mp3",
			ExpectedTrack: Track{
				AlbumArtist: "arcade fire",
				Album:       "funeral",
				TrackNumber: 7,
				Artist:      "arcade fire",
				Title:       "wake up",
			},
			ExpectMatch: true,
		},
		{
			Location:      "fred",
			ExpectedTrack: Track{},
			ExpectMatch:   false,
		},
	}

	ds := &DiskStorage{
		ID:       uuid.NewString(),
		BasePath: "/__DOES_NOT__EXIST__/basepath",
	}
	err := ds.setRegexps(mp3Regexps)
	assert.NoError(t, err, "failed to compile regexps %+v", mp3Regexps)

	for i, testCase := range testCases {
		match := ds.matchLocation(testCase.Location)
		if testCase.ExpectMatch {
			require.NotNil(t, match, "[i=%d] %s got=%+v", i, testCase.Location, match)
			assert.Equal(t, &testCase.ExpectedTrack, match, "[i=%d] %s", i, testCase.Location)
		} else {
			assert.Nil(t, match)
		}
	}
}

func TestAnnotateTrack(t *testing.T) {
	ds := &DiskStorage{
		ID:       uuid.NewString(),
		BasePath: "/__DOES_NOT__EXIST__/basepath",
	}
	err := ds.setRegexps(mp3Regexps)
	assert.NoError(t, err, "failed to compile regexps %+v", mp3Regexps)

	// Annotate track that has tags, except the album artist.
	// This should be able to pull the album artist from the filename.
	t.Run("UseTags", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "artist-path", "not-this-album", "track1.blarg"),
			Tags: Tags{
				Title:  "title",
				Artist: "artist",
				Album:  "album",
			},
		}

		expectedTrack := track
		expectedTrack.Title = track.Tags.Title
		expectedTrack.Name = track.Tags.Artist + " :: " + track.Tags.Title // artist different from album artist
		expectedTrack.Artist = track.Tags.Artist
		expectedTrack.Album = track.Tags.Album
		expectedTrack.AlbumArtist = "artist-path" // from Location, because no AlbumArtist tag
		expectedTrack.PlaylistLocation = "tags:" + filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	// Annotate track that has tags, except the album artist.
	// There is no artist or album information in the filename,
	// so the album artist should default to the artist.
	t.Run("UseTagsFlatDirectory", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "track1.blarg"),
			Tags: Tags{
				Title:  "title",
				Artist: "artist",
				Album:  "album",
			},
		}

		expectedTrack := track
		expectedTrack.Title = track.Tags.Title
		expectedTrack.Name = expectedTrack.Title
		expectedTrack.Artist = track.Tags.Artist
		expectedTrack.Album = track.Tags.Album
		expectedTrack.AlbumArtist = track.Tags.Artist
		expectedTrack.PlaylistLocation = "tags:" + filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	// Annotate track that has tags, except the album artist.
	// There is not enough information in the filename,
	// so the album artist should default to the artist.
	t.Run("UseTagsSingleDirectory", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "sudir", "track1.blarg"),
			Tags: Tags{
				Title:  "title",
				Artist: "artist",
				Album:  "album",
			},
		}

		expectedTrack := track
		expectedTrack.Title = track.Tags.Title
		expectedTrack.Name = expectedTrack.Title
		expectedTrack.Artist = track.Tags.Artist
		expectedTrack.Album = track.Tags.Album
		expectedTrack.AlbumArtist = track.Tags.Artist
		expectedTrack.PlaylistLocation = "tags:" + filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	// Verify that the CDDB disc ID is used as the unique ID in the playlist location.
	t.Run("UseTagsAndCDDBID", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "artist-path", "not-this-album", "track1.blarg"),
			Tags: Tags{
				Title:   "title",
				Artist:  "artist",
				Album:   "album",
				AlbumId: "f00dd00d",
			},
		}

		expectedTrack := track
		expectedTrack.Title = track.Tags.Title
		expectedTrack.Name = track.Tags.Artist + " :: " + track.Tags.Title // artist different from album artist
		expectedTrack.Artist = track.Tags.Artist
		expectedTrack.Album = track.Tags.Album
		expectedTrack.AlbumId = track.Tags.AlbumId
		expectedTrack.AlbumArtist = "artist-path" // from Location, because no AlbumArtist tag
		expectedTrack.PlaylistLocation = "tags:" + filepath.Join(ds.BasePath, expectedTrack.AlbumId, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	t.Run("UseTagsAndAlbumArtist", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "not-this-artist", "not-this-album", "track1.blarg"),
			Tags: Tags{
				Title:       "title",
				Artist:      "artist",
				Album:       "album",
				AlbumArtist: "album-artist",
			},
		}

		expectedTrack := track
		expectedTrack.Title = track.Tags.Title
		expectedTrack.Name = track.Tags.Artist + " :: " + track.Tags.Title // artist different from album artist
		expectedTrack.Artist = track.Tags.Artist
		expectedTrack.Album = track.Tags.Album
		expectedTrack.AlbumArtist = track.Tags.AlbumArtist
		expectedTrack.PlaylistLocation = "tags:" + filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	// Verify that the CDDB disc ID is used in the playlist location.
	// It's more unique than the album artist.
	t.Run("UseTagsAndCDDBIDAndAlbumArtist", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "not-this-artist", "not-this-album", "track1.blarg"),
			Tags: Tags{
				Title:       "title",
				Artist:      "artist",
				Album:       "album",
				AlbumId:     "f00dd00d",
				AlbumArtist: "album-artist",
			},
		}

		expectedTrack := track
		expectedTrack.Title = track.Tags.Title
		expectedTrack.Name = track.Tags.Artist + " :: " + track.Tags.Title // artist different from album artist
		expectedTrack.Artist = track.Tags.Artist
		expectedTrack.Album = track.Tags.Album
		expectedTrack.AlbumId = track.Tags.AlbumId
		expectedTrack.AlbumArtist = track.Tags.AlbumArtist
		expectedTrack.PlaylistLocation = "tags:" + filepath.Join(ds.BasePath, expectedTrack.AlbumId, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	t.Run("UseRegex", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "adam f - colours (08) - circles.mp3"),
		}

		expectedTrack := track
		expectedTrack.Title = "circles"
		expectedTrack.Name = expectedTrack.Title
		//expectedTrack.TrackNumber = 8 // TODO: track number
		expectedTrack.Artist = "adam f"
		expectedTrack.Album = "colours"
		expectedTrack.AlbumArtist = expectedTrack.Artist
		expectedTrack.PlaylistLocation = "regex:" + filepath.Join(ds.BasePath, expectedTrack.Artist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	t.Run("UseRegexAlbumArtist", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "botchit & scarper - botchit breaks (disc 01) (07) - full moon scientist - we demand a shrubbery (scissorkicks mix).mp3"),
		}

		expectedTrack := track
		expectedTrack.Title = "we demand a shrubbery (scissorkicks mix)"
		expectedTrack.Artist = "full moon scientist"
		expectedTrack.Name = expectedTrack.Artist + " :: " + expectedTrack.Title // artist different from album artist
		expectedTrack.Album = "botchit breaks (disc 01)"
		expectedTrack.AlbumArtist = "botchit & scarper"
		expectedTrack.PlaylistLocation = "regex:" + filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)

		// accidental multiple artists, although they're the same:
		// arcade fire - funeral (07) - arcade fire - wake up
	})

	t.Run("UseRegexAccidentalAlbumArtist", func(t *testing.T) {
		// accidental multiple artists, although they're the same:
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "arcade fire - funeral (07) - arcade fire - wake up.mp3"),
		}

		expectedTrack := track
		expectedTrack.Title = "wake up"
		expectedTrack.Name = expectedTrack.Title // because album artist and artist are actually the same
		expectedTrack.Artist = "arcade fire"
		expectedTrack.Album = "funeral"
		expectedTrack.AlbumArtist = "arcade fire"
		expectedTrack.PlaylistLocation = "regex:" + filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})

	t.Run("UseFilename", func(t *testing.T) {
		track := Track{
			ID:       uuid.NewString(),
			Location: filepath.Join(ds.BasePath, "artist", "album", "track1.blarg"),
		}

		expectedTrack := track
		expectedTrack.Title = "track1"
		expectedTrack.Name = expectedTrack.Title
		expectedTrack.Artist = "artist"
		expectedTrack.Album = "album"
		expectedTrack.AlbumArtist = expectedTrack.Artist
		expectedTrack.PlaylistLocation = filepath.Join(ds.BasePath, expectedTrack.AlbumArtist, expectedTrack.Album)

		resultTrack := track
		ds.annotateTrack(&resultTrack)
		assert.Equal(t, expectedTrack, resultTrack)
	})
}
