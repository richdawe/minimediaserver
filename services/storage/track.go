package storage

type Track struct {
	ID       string // ID unique within this storage service
	Location string // Location within storage service (e.g.: filename, URL)
	MIMEType string // MIME type for data, see https://www.iana.org/assignments/media-types/media-types.xhtml#audio
	DataLen  int64  // Size of track data

	Tags Tags // Tags (if any), from track data or elsewhere (e.g.: DB)

	// The following fields may be computed.
	Name        string // Textual description
	Title       string
	Artist      string
	Album       string
	Genre       string // May be empty
	TrackNumber int    // 0 means unknown.
}
