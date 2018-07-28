package xversion

import (
	"fmt"
	"time"

	"strings"

	"github.com/pkg/errors"
	"github.com/unchainio/pkg/xexec"
)

type Version struct {
	Major int
	Minor int
	Patch int

	BuildDate string
	GitHash   string
}

func (v *Version) String() string {
	return fmt.Sprintf(
		` version: %d.%d.%d
 build date: %s
 git hash: %s`, v.Major, v.Minor, v.Patch, v.BuildDate, v.GitHash)
}

func (v *Version) ShortString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Bump() error {
	var err error
	v.Patch += 1
	v.BuildDate = time.Now().Format(time.UnixDate)

	gitHashBytes, err := xexec.Run("git rev-parse --short HEAD 2>/dev/null")

	if err != nil {
		return errors.Wrap(err, "failed to parse git hash")
	}

	v.GitHash = strings.Trim(string(gitHashBytes), "\n")

	return nil
}
