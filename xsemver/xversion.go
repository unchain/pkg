package xsemver

import (
	"fmt"
	"time"

	"github.com/coreos/go-semver/semver"
)

type Version struct {
	sv *semver.Version

	BuildDate string
}

type Versions []*Versions

func NewVersion(v string) (*Version, error) {
	sv, err := semver.NewVersion(v)

	if err != nil {
		return nil, err
	}

	vv := &Version{
		sv:        sv,
		BuildDate: time.Now().Format(time.UnixDate),
	}

	return vv, nil
}

func (v *Version) String() string {
	return fmt.Sprintf(
		` version: %d.%d.%d
 build date: %s`, v.sv.Major, v.sv.Minor, v.sv.Patch, v.BuildDate)
}

func (v *Version) ShortString() string {
	return fmt.Sprintf("%d.%d.%d", v.sv.Major, v.sv.Minor, v.sv.Patch)
}

func (v *Version) BumpMajor() {
	v.sv.BumpMajor()

	v.BuildDate = time.Now().Format(time.UnixDate)
}

func (v *Version) BumpMinor() {
	v.sv.BumpMinor()

	v.BuildDate = time.Now().Format(time.UnixDate)
}

func (v *Version) BumpPatch() {
	v.sv.BumpPatch()

	v.BuildDate = time.Now().Format(time.UnixDate)
}
