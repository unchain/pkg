package xtest

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadBytes is a test helper that returns the bytes of a file with path relative to testdata/
func LoadBytes(t *testing.T, path string) []byte {
	bytes, err := ioutil.ReadFile(filepath.Join("testdata", path))
	require.NoError(t, err)

	return bytes
}

// LoadBytes is a test helper that returns the bytes of a file with path relative to testdata/
func WriteBytes(t *testing.T, path string, bytes []byte) {
	err := ioutil.WriteFile(filepath.Join("testdata", path), bytes, 0644)
	require.NoError(t, err)
}
