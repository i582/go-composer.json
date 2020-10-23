package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int64
	Minor int64
	Micro int64

	IsDev   bool
	IsPatch bool
	IsAlpha bool
	IsBeta  bool
	IsRC    bool
}

func NewVersion(val string) (*Version, error) {
	var version = &Version{}

	if val == "" {
		return nil, fmt.Errorf("version is empty")
	}

	if len(val) < 5 {
		return nil, fmt.Errorf("version must be in the format [v]X.Y.Z[-suffix]")
	}

	val = strings.TrimPrefix(val, "v")
	vals := strings.Split(val, "-")
	if len(vals) > 2 {
		return nil, fmt.Errorf("version must be in the format [v]X.Y.Z[-suffix]")
	}

	if len(vals) == 2 {
		switch vals[1] {
		case "dev":
			version.IsDev = true
		case "patch", "p":
			version.IsPatch = true
		case "alpha", "a":
			version.IsAlpha = true
		case "beta", "b":
			version.IsBeta = true
		case "RC":
			version.IsRC = true
		default:
			return nil, fmt.Errorf("unknown version suffix '%s'", vals[1])
		}
	}

	val = vals[0]

	vals = strings.Split(val, ".")

	if len(vals) != 3 {
		return nil, fmt.Errorf("version must be in the format [v]X.Y.Z[-suffix]")
	}

	major, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("part 1 ('%s') of the version must be a number", vals[0])
	}

	minor, err := strconv.ParseInt(vals[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("part 2 ('%s') of the version must be a number", vals[1])
	}

	micro, err := strconv.ParseInt(vals[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("part 3 ('%s') of the version must be a number", vals[2])
	}

	version.Major = major
	version.Minor = minor
	version.Micro = micro

	return version, nil
}

func (v *Version) HasPrefix() bool {
	return v.IsDev != false || v.IsPatch != false || v.IsAlpha != false || v.IsBeta != false || v.IsRC != false
}
