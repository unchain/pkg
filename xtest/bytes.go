package xtest

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// LoadBytes is a test helper that returns the bytes of a file with path relative to testdata/
func LoadBytes(t *testing.T, path string) []byte {
	bytes, err := ioutil.ReadFile(filepath.Join("testdata", path))

	if err != nil {
		t.Fatal(err)
	}

	return bytes
}
