package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/google/uuid"
)

type DiskStorage struct {
	ID       string
	BasePath string

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

// TODO: test coverage for findPlaylist*()
func findPlaylistLocation(location string) string {
	return filepath.Dir(location)
}

func findPlaylistArtistAlbum(location string) (string, string) {
	playlistLocation := findPlaylistLocation(location)

	album := filepath.Base(playlistLocation)
	artist := filepath.Base(filepath.Dir(playlistLocation))
	return artist, album
}

func (ds *DiskStorage) buildTracks() (map[string]Track, map[string]Playlist, error) {
	tracksByID := make(map[string]Track, 0)
	playlistsByID := make(map[string]Playlist, 0)
	playlistsByLocation := make(map[string]string, 0) // value is playlist ID

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
		name := tags.Title
		if name == "" {
			name = d.Name()
		}

		track := Track{
			Name:     name,
			Tags:     tags,
			ID:       trackUUID,
			Location: location,
			MIMEType: mimeType,
			DataLen:  fileinfo.Size(),
		}
		tracksByID[track.ID] = track

		// TODO: This builds playlists based on albums having their own directory.
		// Other music collections (e.g.: my MP3s) use a flat format
		// with everything encoded in one filename. Could use with the playlist building
		// strategy being pluggable.
		playlistLocation := findPlaylistLocation(location)
		playlistUUID := locationToUUIDString(playlistLocation)

		playlistID, ok := playlistsByLocation[playlistLocation]
		if !ok {
			artist, album := findPlaylistArtistAlbum(location)
			playlistName := fmt.Sprintf("%s :: %s", artist, album)

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
		return nil
	})

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

	// Find and sort the playlist IDs based on the location
	// of the playlist. Then build list of sorted playlists
	// using the sorted list of playlist IDs.
	//
	// TODO: This might need a pluggable strategy for music
	// that is not split out into a directory per album
	// (see also the comment in the building process).
	playlistIDs := make([]string, 0, 1)
	for _, playlist := range playlistsByID {
		playlistIDs = append(playlistIDs, playlist.ID)
	}
	sort.Slice(playlistIDs, func(i int, j int) bool {
		playlistI := playlistsByID[playlistIDs[i]]
		playlistJ := playlistsByID[playlistIDs[j]]
		return playlistI.Location < playlistJ.Location
	})

	playlists := make([]Playlist, 0, 1)
	for _, playlistID := range playlistIDs {
		playlists = append(playlists, playlistsByID[playlistID])
	}

	// TODO: sort tracks in playlist too

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

func NewDiskStorage(path string) (*DiskStorage, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	return &DiskStorage{
		ID:       uuid.New().String(),
		BasePath: path,
	}, nil
}
