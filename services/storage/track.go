package storage

type Track struct {
	ID       string // ID unique within this storage service
	Location string // Location within storage service (e.g.: filename, URL)

	Name     string // Textual description
	Tags     Tags   // Tags (if any), from track data or elsewhere (e.g.: DB)
	MIMEType string // MIME type for data, see https://www.iana.org/assignments/media-types/media-types.xhtml#audio
	DataLen  int64  // Size of track data
}
