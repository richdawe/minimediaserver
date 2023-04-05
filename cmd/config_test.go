package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	// Default configuration (no config from files)
	config, err := loadConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	t.Run("Default config", func(t *testing.T) {
		assert.Equal(t, "127.0.0.1:1323", config.Addr)
		assert.Len(t, config.StorageServices, 0)

		catalogService, err := buildCatalog(config)
		require.NoError(t, err)
		require.NotNil(t, catalogService)
	})

	t.Run("Add some storage services", func(t *testing.T) {
		config.StorageServices = append(config.StorageServices, StorageServiceConfig{
			Type: "nullStorage", Path: "/dev/null",
		})
		config.StorageServices = append(config.StorageServices, StorageServiceConfig{
			Type: "diskStorage", Path: "",
		})

		catalogService, err := buildCatalog(config)
		require.NoError(t, err)
		require.NotNil(t, catalogService)
	})

	t.Run("Invalid storage service name", func(t *testing.T) {
		config.StorageServices = append(config.StorageServices, StorageServiceConfig{
			Type: "MEGASTORAGE!!!!1", Path: "/dev/urandom",
		})

		_, err := buildCatalog(config)
		require.Error(t, err)
	})

	t.Run("Invalid file path for disk storage", func(t *testing.T) {
		config.StorageServices = append(config.StorageServices, StorageServiceConfig{
			Type: "diskStorage", Path: "/dev/__DOES_NOT_EXIST__",
		})

		_, err := buildCatalog(config)
		require.Error(t, err)
	})
}
