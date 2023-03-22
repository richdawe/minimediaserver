package storage

type Track struct {
	Name     string // Textual description
	ID       string // ID unique within this storage service
	Location string // Location within storage service (e.g.: filename, URL)
	MIMEType string // MIME type for data, see https://www.iana.org/assignments/media-types/media-types.xhtml#audio
}
