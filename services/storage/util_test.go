package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	t.Run("getMIMEType", func(t *testing.T) {
		testCases := []struct {
			Filename string
			Expected string
		}{
			{"foo.mp3", "audio/mp3"},
			{"foo.mp3s", "application/binary"},
			{"foo.m4a", "audio/mp4"},
			{"foo.ogg", "audio/ogg"},
			{"foo.flac", "audio/flac"},
			{"foo.txt", "application/binary"},
		}

		for _, testCase := range testCases {
			assert.Equal(t, testCase.Expected, getMIMEType(testCase.Filename), testCase.Filename)
		}
	})

	t.Run("ignoreMIMEType", func(t *testing.T) {
		testCases := []struct {
			MIMEType string
			Expected bool
		}{
			{"audio/mp3", false},
			{"audio/mp4", false},
			{"audio/ogg", false},
			{"audio/flac", false},
			{"application/binary", true},
		}

		for _, testCase := range testCases {
			assert.Equal(t, testCase.Expected, ignoreMIMEType(testCase.MIMEType), testCase.MIMEType)
		}
	})
}
