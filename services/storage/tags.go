package storage

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/jfreymuth/oggvorbis"
)

type Tags struct {
	Title  string
	Artist string
	Album  string
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
		tags.Album = artist
	}
	if album, ok := cm[flacvorbis.FIELD_ALBUM]; ok {
		tags.Album = album
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
