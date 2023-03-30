package storage

import (
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
