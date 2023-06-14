package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richdawe/minimediaserver/services/catalog"
	"github.com/richdawe/minimediaserver/services/storage"
)

func TestEndpoints(t *testing.T) {
	var config Config

	catalogService, err := catalog.NewBasicCatalog()
	require.NoError(t, err)

	// TODO: need a config file for configuring storage backends
	nullStorage, err := storage.NewNullStorage()
	require.NoError(t, err)
	err = catalogService.AddStorage(nullStorage)
	require.NoError(t, err)

	e, err := setupEndpoints(config, catalogService)
	require.NoError(t, err)
	require.NotNil(t, e) // TODO: remove when something more interesting is happening

	// TODO: exercise more HTTP endpoints

	t.Run("StaticEndpoints", func(t *testing.T) {
		// TODO: starting point at https://echo.labstack.com/guide/testing/
	})
}
