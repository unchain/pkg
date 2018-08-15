package xmerge_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/unchainio/pkg/iferr"
	"github.com/unchainio/pkg/xmerge"
)

type Input1 struct {
	A int
	B int
}

type Input2 struct {
	A int
}

func TestName(t *testing.T) {
	i1 := &Input1{
		1, 2,
	}

	i2 := &Input1{
		A: 7,
	}

	err := xmerge.Merge(i1, i2)

	iferr.Exit(err)

	spew.Dump(i1, i2)
}
