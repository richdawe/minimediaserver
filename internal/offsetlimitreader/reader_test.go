package offsetlimitreader_test

import (
	"io"
	"strings"
	"testing"

	"github.com/richdawe/minimediaserver/internal/offsetlimitreader"
	"github.com/stretchr/testify/assert"
)

// TODO: test string >> buffer size
func TestSmall(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Start  int64
		Length int64

		ExpectEOF   bool
		ExpectedLen int
		Expected    string
	}{
		{
			Name:      "Empty string",
			Input:     "",
			Start:     0,
			Length:    0,
			ExpectEOF: true,
		},
		{
			Name:      "Empty string, length longer than string",
			Input:     "",
			Start:     0,
			Length:    1,
			ExpectEOF: true,
		},
		{
			Name:        "Single char",
			Input:       "A",
			Start:       0,
			Length:      1,
			ExpectedLen: 1,
			Expected:    "A",
		},
		{
			Name:        "Two chars",
			Input:       "AB",
			Start:       0,
			Length:      2,
			ExpectedLen: 2,
			Expected:    "AB",
		},
		{
			Name:        "Limit to last char",
			Input:       "AB",
			Start:       1,
			Length:      1,
			ExpectedLen: 1,
			Expected:    "B",
		},
		{
			Name:        "Limit to middle char",
			Input:       "ABC",
			Start:       1,
			Length:      1,
			ExpectedLen: 1,
			Expected:    "B",
		},
		{
			Name:        "Limit to first char",
			Input:       "ABC",
			Start:       0,
			Length:      1,
			ExpectedLen: 1,
			Expected:    "A",
		},
		{
			Name:        "Middle of a longer string",
			Input:       "123_456_789",
			Start:       4,
			Length:      3,
			ExpectedLen: 3,
			Expected:    "456",
		},
		{
			Name:      "Limit at end",
			Input:     "hello",
			Start:     5,
			Length:    0,
			ExpectEOF: true,
		},
		{
			Name:      "Limit at end, length off end of string",
			Input:     "hello",
			Start:     5,
			Length:    1,
			ExpectEOF: true,
		},
	}

	for _, testCase := range testCases {
		r := strings.NewReader(testCase.Input)
		olr := offsetlimitreader.New(r, testCase.Start, testCase.Length)

		buf := make([]byte, 128)
		n, err := olr.Read(buf)

		assert.Equal(t, testCase.ExpectedLen, n, testCase.Name)
		if testCase.ExpectEOF {
			assert.Equal(t, io.EOF, err, testCase.Name)
			continue
		}
		assert.Nil(t, err, testCase.Name)

		read := string(buf[:n])
		if n == testCase.ExpectedLen {
			assert.Equal(t, testCase.Expected, read, testCase.Name)
		}
	}
}
