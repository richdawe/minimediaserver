package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtil(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		uuid1 := locationToUUIDString("file1")
		uuid2 := locationToUUIDString("file2")
		require.NotEmpty(t, uuid1)
		require.NotEmpty(t, uuid2)
		require.NotEqual(t, uuid1, uuid2)
	})
}
