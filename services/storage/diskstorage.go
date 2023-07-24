package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type DiskStorage struct {
	ID       string
	BasePath string
	Regexps  []string

	compiledRegexps []*regexp.Regexp

	tracksByID    map[string]Track
	playlistsByID map[string]Playlist

	sortedTracks    []Track
	sortedPlaylists []Playlist
}

func (ds *DiskStorage) GetID() string {
	return ds.ID
}

// Find the tracks in this storage, and return the tracks
// in a stable order.
func (ds *DiskStorage) FindTracks() ([]Track, []Playlist, error) {
	var err error

	if ds.tracksByID == nil {
		ds.tracksByID, ds.playlistsByID, err = ds.buildTracks()
	}
	if err != nil {
		return nil, nil, err
	}
	if ds.sortedTracks == nil {
		ds.sortedTracks, ds.sortedPlaylists = ds.buildSortedTracks()
	}
	return ds.sortedTracks, ds.sortedPlaylists, nil
}

func buildPlaylists(tracksByID map[string]Track) (map[string]Playlist, error) {
	playlistsByID := make(map[string]Playlist, 0)
	playlistsByLocation := make(map[string]string, 0) // value is playlist ID

	for _, track := range tracksByID {
		playlistLocation := track.PlaylistLocation
		playlistUUID := locationToUUIDString(playlistLocation)

		playlistID, ok := playlistsByLocation[playlistLocation]
		if !ok {
			// TODO: This isn't quite right in the case where there is an album
			// with multiple artists, but the AlbumArtist defaults to the artist.
			// We need some heuristics to group tracks based on location / album ID / album name
			// and then generate "Various artists" if each track in the playlist
			// has a different artist name.
			//
			// This will probably require track.AlbumArtist to be empty,
			// if it's not defined in tags, so that we know when to use heuristics.
			//
			// Note: This has been partially implemented with use of CDDB ID and/or track location.
			playlistName := fmt.Sprintf("%s :: %s", track.AlbumArtist, track.Album)

			playlist := Playlist{
				Name:     playlistName,
				ID:       playlistUUID,
				Location: playlistLocation,
				Tracks:   make([]Track, 0, 1),
			}

			playlistID = playlist.ID
			playlistsByLocation[playlistLocation] = playlistID
			playlistsByID[playlistID] = playlist
		}

		playlist := playlistsByID[playlistID]
		playlist.Tracks = append(playlist.Tracks, track)
		playlistsByID[playlistID] = playlist
	}

	// Sort the tracks in a playlist by the filename (location).
	// This is a proxy for their position in the playlist (until we use tags).
	for _, playlist := range playlistsByID {
		sort.Slice(playlist.Tracks, func(i int, j int) bool {
			return playlist.Tracks[i].Location < playlist.Tracks[j].Location
		})
	}

	return playlistsByID, nil
}

// Determine whether the track artist and album artist are the same or similar,
// using some normalizations and ignoring case.
func isTrackByAlbumArtist(trackArtist string, albumArtist string) bool {
	// Replace underscore and dash with space,
	// since those may commonly be used as replacements for space
	// in filenames.
	//
	// Also to consider, to improve matching:
	// - removing apostrophes
	//
	replacer := func(in string) string {
		out := strings.ReplaceAll(in, "_", " ")
		out = strings.ReplaceAll(out, "-", " ")
		return out
	}

	a := replacer(trackArtist)
	b := replacer(albumArtist)
	return strings.EqualFold(a, b)
}

func removeFileExtension(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		idx = len(filename)
	}
	return filename[:idx]
}

// Use any regular expressions for this storage to match track artist,
// etc. from the filename.
func (ds *DiskStorage) matchLocation(location string) *Track {
	if len(ds.compiledRegexps) == 0 {
		return nil
	}

	var t *Track
	filename := filepath.Base(location)
	filename = removeFileExtension(filename)

	for _, c := range ds.compiledRegexps {
		matches := c.FindStringSubmatch(filename)
		if matches == nil {
			continue
		}

		t = &Track{}
		if idx := c.SubexpIndex("albumartist"); idx != -1 {
			t.AlbumArtist = strings.Trim(matches[idx], " ")
		}
		if idx := c.SubexpIndex("album"); idx != -1 {
			t.Album = strings.Trim(matches[idx], " ")
		}
		if idx := c.SubexpIndex("trackno"); idx != -1 {
			t.TrackNumber, _ = strconv.Atoi(matches[idx])
		}
		if idx := c.SubexpIndex("artist"); idx != -1 {
			t.Artist = strings.Trim(matches[idx], " ")
		}
		if idx := c.SubexpIndex("title"); idx != -1 {
			t.Title = strings.Trim(matches[idx], " ")
		}
		break
	}

	return t
}

// TODO: Probably would be cleaner to have this build the track completely
// given some input data?
func (ds *DiskStorage) annotateTrack(track *Track) {
	var artist, album, albumArtist, albumId, title string

	// Default playlist location is the directory containing the file.
	// This may be overridden below.
	playlistLocation := filepath.Dir(track.Location)
	// TODO: need a default playlist name here too (e.g.: for "mp3" directory)
	// in case the strategy doesn't work.

	// Strategy 1: Use tags to determine artist, album, etc. information.
	if track.Tags.Artist != "" && track.Tags.Album != "" && track.Tags.Title != "" {
		artist = track.Tags.Artist
		album = track.Tags.Album
		albumArtist = track.Tags.AlbumArtist
		albumId = track.Tags.AlbumId
		// TODO: track number
		title = track.Tags.Title

		// TODO: track number, and use that to position in playlists.

		// Heuristic: If the album artist wasn't determined by tags or regex,
		// use the directory name. But only when the filename is like
		// /basepath/artist/album/filename.flac ,
		// /basepath/subdir/artist/album/filename.flac , or similar.
		if albumArtist == "" {
			trackDir := filepath.Dir(track.Location)
			albumDir := filepath.Dir(trackDir)
			if trackDir != ds.BasePath && albumDir != ds.BasePath {
				albumArtist = filepath.Base(albumDir)
			}
		}

		// Use the most defined tag we can for the playlist path.
		artistPath := artist
		if albumArtist != "" {
			artistPath = albumArtist
		}
		if albumId != "" {
			artistPath = albumId
		}

		playlistLocation = "tags:" + filepath.Join(ds.BasePath, artistPath, album)
	}

	// Strategy 2: Regular expression matching (if enabled for this storage instance).
	if artist == "" && album == "" && title == "" {
		if t := ds.matchLocation(track.Location); t != nil {
			album = t.Album
			artist = t.Artist
			title = t.Title
			// TODO: track number
			if t.AlbumArtist != "" {
				albumArtist = t.AlbumArtist
			}

			// Use the most defined tag we can for the playlist path.
			artistPath := artist
			if albumArtist != "" {
				artistPath = albumArtist
			}

			playlistLocation = "regex:" + filepath.Join(ds.BasePath, artistPath, album)
		}
	}

	// Strategy 3: Parse from directory and filename (heuristic)
	if artist == "" && album == "" && title == "" {
		trackDir := filepath.Dir(track.Location)
		album = filepath.Base(trackDir)
		artist = filepath.Base(filepath.Dir(trackDir))

		filename := filepath.Base(track.Location)
		title = removeFileExtension(filename)
	}

	// Finally, annotate the track.
	if albumArtist == "" {
		albumArtist = artist
	}

	track.Artist = artist
	track.Album = album
	track.AlbumArtist = albumArtist
	track.AlbumId = albumId
	track.Title = title

	// Determine if the track name should include the track's artist,
	// for multi-artist albums.
	if isTrackByAlbumArtist(track.Artist, track.AlbumArtist) {
		track.Name = track.Title
	} else {
		track.Name = fmt.Sprintf("%s :: %s", track.Artist, track.Title)
	}

	track.PlaylistLocation = playlistLocation
}

func (ds *DiskStorage) buildTracks() (map[string]Track, map[string]Playlist, error) {
	tracksByID := make(map[string]Track, 0)

	fileSystem := os.DirFS(ds.BasePath)

	walkErr := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		location := filepath.Join(ds.BasePath, path)
		trackUUID := locationToUUIDString(location)

		// Ignore some unknown MIME types
		mimeType := getMIMEType(d.Name())
		if ignoreMIMEType(mimeType) {
			fmt.Printf("Ignoring file due to MIME type: %s\n", location)
			return nil
		}

		fileinfo, err := os.Stat(location)
		if err != nil {
			return err
		}

		// TODO: move tags handling into common code for storage engines
		r, err := os.Open(location)
		if err != nil {
			return err
		}
		tags, err := readTags(r, mimeType)
		if err != nil {
			return err
		}

		track := Track{
			ID:       trackUUID,
			Location: location,
			MIMEType: mimeType,
			DataLen:  fileinfo.Size(),
			Tags:     tags,
		}
		ds.annotateTrack(&track)
		tracksByID[track.ID] = track

		return nil
	})

	playlistsByID, err := buildPlaylists(tracksByID)
	if err != nil {
		return nil, nil, err
	}

	return tracksByID, playlistsByID, walkErr
}

func (ds *DiskStorage) buildSortedTracks() ([]Track, []Playlist) {
	tracksByID := ds.tracksByID
	playlistsByID := ds.playlistsByID

	if tracksByID == nil || playlistsByID == nil {
		return make([]Track, 0), make([]Playlist, 0)
	}

	// Find and sort the track IDs based on the location
	// of the track. Then build list of sorted tracks
	// using the sorted list of track IDs.
	trackIDs := make([]string, 0, 1)
	for _, track := range tracksByID {
		trackIDs = append(trackIDs, track.ID)
	}
	sort.Slice(trackIDs, func(i int, j int) bool {
		trackI := tracksByID[trackIDs[i]]
		trackJ := tracksByID[trackIDs[j]]
		return trackI.Location < trackJ.Location
	})

	tracks := make([]Track, 0, 1)
	for _, trackID := range trackIDs {
		tracks = append(tracks, tracksByID[trackID])
	}

	// Find and sort the playlist IDs based on the name
	// of the playlist. Then build list of sorted playlists
	// using the sorted list of playlist IDs.
	playlistIDs := make([]string, 0, 1)
	for _, playlist := range playlistsByID {
		playlistIDs = append(playlistIDs, playlist.ID)
	}
	sort.Slice(playlistIDs, func(i int, j int) bool {
		playlistI := playlistsByID[playlistIDs[i]]
		playlistJ := playlistsByID[playlistIDs[j]]
		return playlistI.Name < playlistJ.Name
	})

	playlists := make([]Playlist, 0, 1)
	for _, playlistID := range playlistIDs {
		playlists = append(playlists, playlistsByID[playlistID])
	}

	// TODO: sort tracks in playlist too, to match
	// track ordering from tags.

	return tracks, playlists
}

func (ds *DiskStorage) ReadTrack(id string) (io.Reader, error) {
	track, ok := ds.tracksByID[id]
	if !ok {
		// TODO: look at standardizing errors
		return nil, errors.New("track not found")
	}

	// TODO: figure out some way to return a reader that reads in the file in chunks
	data, err := os.ReadFile(track.Location)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func (ds *DiskStorage) setRegexps(regexps []string) error {
	ds.Regexps = make([]string, 0)
	ds.compiledRegexps = make([]*regexp.Regexp, 0)

	for _, r := range regexps {
		c, err := regexp.Compile(r)
		if err != nil {
			// TODO: perhaps this should be a non-fatal error during start-up -
			// log a warning but carry on?
			return err
		}
		ds.Regexps = append(ds.Regexps, r)
		ds.compiledRegexps = append(ds.compiledRegexps, c)
	}

	return nil
}

func NewDiskStorage(path string, regexps []string) (*DiskStorage, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	ds := &DiskStorage{
		ID:       uuid.NewString(),
		BasePath: path,
	}
	err = ds.setRegexps(regexps)
	if err != nil {
		return nil, err
	}
	return ds, nil
}
