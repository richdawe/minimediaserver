package storage

import (
	"strings"

	"github.com/google/uuid"
)

var secret string = "you'll never guess this, oops"

// locationToUUIDString converts a location into a stable UUID value,
// for use in HTTP paths.
func locationToUUIDString(location string) string {
	data := location + ":" + secret
	u := uuid.NewSHA1(uuid.NameSpaceURL, []byte(data))
	return u.String()
}

func getMIMEType(filename string) string {
	var mimeType string

	// https://github.com/apache/httpd/blob/trunk/docs/conf/mime.types
	filename = strings.ToLower(filename)
	switch {
	case strings.HasSuffix(filename, ".mp3"):
		mimeType = "audio/mp3"
	case strings.HasSuffix(filename, ".m4a"):
		mimeType = "audio/mp4"
	case strings.HasSuffix(filename, ".ogg"):
		mimeType = "audio/ogg"
	case strings.HasSuffix(filename, ".flac"):
		mimeType = "audio/flac"
	}

	if mimeType == "" {
		mimeType = "application/binary"
	}
	return mimeType
}

func ignoreMIMEType(mimeType string) bool {
	switch mimeType {
	case "application/binary":
		return true
	}
	return false
}
