package catalog

import (
	"github.com/richdawe/minimediaserver/services/storage"
)

type Track struct {
	StorageService storage.StorageService

	ID   string
	Name string
}
