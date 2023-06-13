package storage

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/jfreymuth/oggvorbis"
)

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

// TODO: Unify processing for ogg and flac vorbis comments
func readOggTags(r io.Reader) (Tags, error) {
	var tags Tags

	oggfile, err := oggvorbis.NewReader(r)
	if err != nil {
		return Tags{}, err
	}

	cm := make(map[string]string)

	comments := oggfile.CommentHeader().Comments
	for _, v := range comments {
		p := strings.SplitN(v, "=", 2)
		if len(p) != 2 {
			continue
		}
		ckey := strings.ToUpper(p[0])
		cm[ckey] = p[1]
	}

	if title, ok := cm[flacvorbis.FIELD_TITLE]; ok {
		tags.Title = title
	}
	if artist, ok := cm[flacvorbis.FIELD_ARTIST]; ok {
		tags.Artist = artist
	}
	if album, ok := cm[flacvorbis.FIELD_ALBUM]; ok {
		tags.Album = album
	}
	if genre, ok := cm[flacvorbis.FIELD_GENRE]; ok {
		tags.Genre = genre
	}
	if trackNumber, ok := cm[flacvorbis.FIELD_TRACKNUMBER]; ok {
		n, err := strconv.Atoi(trackNumber)
		if err == nil {
			tags.TrackNumber = n
		}
	}

	return tags, nil
}

// getFlacComment returns the first comment value
func getFlacComment(cmt *flacvorbis.MetaDataBlockVorbisComment, name string) (string, bool) {
	results, _ := cmt.Get(name)
	if len(results) == 0 {
		return "", false
	}
	return results[0], true
}

func readFlacTags(r io.Reader) (Tags, error) {
	var tags Tags

	flacfile, err := flac.ParseMetadata(r)
	if err != nil {
		return Tags{}, err
	}

	for _, meta := range flacfile.Meta {
		if meta.Type == flac.VorbisComment {
			cmt, err := flacvorbis.ParseFromMetaDataBlock(*meta)
			if err != nil {
				return Tags{}, err
			}

			if title, ok := getFlacComment(cmt, flacvorbis.FIELD_TITLE); ok {
				tags.Title = title
			}
			if artist, ok := getFlacComment(cmt, flacvorbis.FIELD_ARTIST); ok {
				tags.Artist = artist
			}
			if album, ok := getFlacComment(cmt, flacvorbis.FIELD_ALBUM); ok {
				tags.Album = album
			}
			if genre, ok := getFlacComment(cmt, flacvorbis.FIELD_GENRE); ok {
				tags.Genre = genre
			}
			if trackNumber, ok := getFlacComment(cmt, flacvorbis.FIELD_TRACKNUMBER); ok {
				n, err := strconv.Atoi(trackNumber)
				if err == nil {
					tags.TrackNumber = n
				}
			}
		}
	}

	return tags, nil
}

func readTags(r io.Reader, mimeType string) (Tags, error) {
	switch mimeType {
	case OggMimeType:
		return readOggTags(r)
	case FlacMimeType:
		return readFlacTags(r)
	case MP3MimeType:
		return Tags{}, nil // TODO: implement ID3 tags
	case MP4MimeType:
		return Tags{}, nil // TODO: implement MP4 tags
	}
	return Tags{}, fmt.Errorf("unable to read tags for MIME type %s", mimeType)
}
