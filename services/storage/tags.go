package storage

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/jfreymuth/oggvorbis"
	id3 "github.com/richdawe/id3-go"
)

// TODO: one day look at replacing this code with https://github.com/dhowden/tag ?

// *** Setting the title, etc. tags for Ogg files:
//
// TITLE comment checked and set using:
// vorbiscomment services/storage/example.ogg
// vorbiscomment -a -t title=ExAmPlE services/storage/example.ogg
//
// See also proposed standard field names at https://www.xiph.org/vorbis/doc/v-comment.html

// *** Setting the tite, etc. for Flac files:
//
// metaflac --set-tag "TITLE=ALBUM1_TRACK2_EXAMPLE" --set-tag "album=album1" --set-tag "ARTIST=the-artist" testdata/services/storage/diskstorage/Music/cds/Artist/Album1/track2-example.flac
// metaflac --list testdata/services/storage/diskstorage/Music/cds/Artist/Album1/track2-example.flac
//
// This uses the vorbis comment format too.

type Tags struct {
	Title       string
	Artist      string
	Album       string
	Genre       string
	TrackNumber int
}

// Convert a vorbis comment list into a map for lookups.
func commentsToMap(comments []string) map[string]string {
	cm := make(map[string]string)

	for _, v := range comments {
		p := strings.SplitN(v, "=", 2)
		if len(p) != 2 {
			continue
		}
		ckey := strings.ToUpper(p[0])
		// Append if key already exists; see e.g.:
		// https://github.com/go-flac/flacvorbis/blob/v0.1.0/vorbis.go#L26
		cm[ckey] += p[1]
	}

	return cm
}

// Get the tags we're interested in from a map of comments.
func getTags(commentsMap map[string]string) (tags Tags) {
	if title, ok := commentsMap[flacvorbis.FIELD_TITLE]; ok {
		tags.Title = title
	}
	if artist, ok := commentsMap[flacvorbis.FIELD_ARTIST]; ok {
		tags.Artist = artist
	}
	if album, ok := commentsMap[flacvorbis.FIELD_ALBUM]; ok {
		tags.Album = album
	}
	if genre, ok := commentsMap[flacvorbis.FIELD_GENRE]; ok {
		tags.Genre = genre
	}
	if trackNumber, ok := commentsMap[flacvorbis.FIELD_TRACKNUMBER]; ok {
		n, err := strconv.Atoi(trackNumber)
		if err == nil {
			tags.TrackNumber = n
		}
	}
	return
}

// Read tags from an OGG file.
func readOggTags(r io.Reader) (Tags, error) {
	oggfile, err := oggvorbis.NewReader(r)
	if err != nil {
		return Tags{}, err
	}

	comments := oggfile.CommentHeader().Comments
	cm := commentsToMap(comments)
	return getTags(cm), nil
}

// Read tags from a FLAC file.
func readFlacTags(r io.Reader) (Tags, error) {
	var tags Tags

	flacfile, err := flac.ParseMetadata(r)
	if err != nil {
		return Tags{}, err
	}

	for _, meta := range flacfile.Meta {
		if meta.Type == flac.VorbisComment {
			vorbisComments, err := flacvorbis.ParseFromMetaDataBlock(*meta)
			if err != nil {
				return Tags{}, err
			}

			cm := commentsToMap(vorbisComments.Comments)
			tags = getTags(cm)
		}
	}

	return tags, nil
}

// Read tags from an MP3 file.
func readMP3Tags(r io.Reader) (Tags, error) {
	// TODO: This is a bit clunky. id3-go doesn't provide
	// an io.Reader or io.ReadSeeker interface. So we try to cast
	// to a type it doesn't accept. This works for the current storage types
	// but may fail for future ones (e.g.: fetching from S3 or a database).
	f, ok := r.(*os.File)
	if !ok {
		return Tags{}, fmt.Errorf("unable to read tags from a non-file")
	}

	file, err := id3.Parse(f)
	if err != nil {
		return Tags{}, err
	}

	return Tags{
		Title:  file.Title(),
		Artist: file.Artist(),
		Album:  file.Album(),
		Genre:  file.Genre(),
	}, nil
}

// Read the tags from a media file.
func readTags(r io.Reader, mimeType string) (Tags, error) {
	switch mimeType {
	case OggMimeType:
		return readOggTags(r)
	case FlacMimeType:
		return readFlacTags(r)
	case MP3MimeType:
		return readMP3Tags(r)
	case MP4MimeType:
		return Tags{}, nil // TODO: implement MP4 tags
	}
	return Tags{}, fmt.Errorf("unable to read tags for MIME type %s", mimeType)
}
