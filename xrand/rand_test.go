package xrand_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/unchainio/pkg/xrand"
)

func TestRandomString(t *testing.T) {
	spew.Dump(xrand.String(100))
}
