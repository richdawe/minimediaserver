package storage

type Track struct {
	ID       string // ID unique within this storage service
	Location string // Location within storage service (e.g.: filename, URL)
	MIMEType string // MIME type for data, see https://www.iana.org/assignments/media-types/media-types.xhtml#audio
	DataLen  int64  // Size of track data

	Tags Tags // Tags (if any), from track data or elsewhere (e.g.: DB)

	// The following fields are computed.
	Name        string // Textual description
	Title       string
	Artist      string
	Album       string
	AlbumArtist string // Will be different to Artist for e.g.: for compilations, or orchestral performances
	AlbumId     string // May be empty
	Genre       string // May be empty
	TrackNumber int    // 0 means unknown.

	PlaylistLocation string // Location for the playlist; may be a virtual URL, like tags:/path or regex:/path
}
