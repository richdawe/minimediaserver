package catalog

type Track struct {
	ID               string
	StorageServiceID string
	Name             string
	MIMEType         string // MIME type for data, see https://www.iana.org/assignments/media-types/media-types.xhtml#audio
}
