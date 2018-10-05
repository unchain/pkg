package xsemver_test

import (
	"fmt"
	"testing"

	"github.com/unchainio/pkg/xsemver"
)

func TestVersion(t *testing.T) {
	version, _ := xsemver.NewVersion("1.2.3-alpha")

	fmt.Print(version)

}
