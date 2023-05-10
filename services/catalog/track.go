package catalog

type Track struct {
	ID               string // Unique ID from storage service
	StorageServiceID string // Storage service's ID

	Name     string
	MIMEType string // MIME type for data, see https://www.iana.org/assignments/media-types/media-types.xhtml#audio
	DataLen  int64  // Size of track data
}
