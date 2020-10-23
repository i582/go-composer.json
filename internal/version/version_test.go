package version

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		Version  string
		Expected Version
		Error    error
	}{
		{
			Version: "2.5.1",
			Expected: Version{
				Major: 2,
				Minor: 5,
				Micro: 1,
			},
		},
		{
			Version: "v0.5.1",
			Expected: Version{
				Major: 0,
				Minor: 5,
				Micro: 1,
			},
		},
		{
			Version: "v0.5.1-dev",
			Expected: Version{
				Major: 0,
				Minor: 5,
				Micro: 1,
				IsDev: true,
			},
		},
		{
			Version: "v0.5.1-patch",
			Expected: Version{
				Major:   0,
				Minor:   5,
				Micro:   1,
				IsPatch: true,
			},
		},
		{
			Version: "v0.5.1-p",
			Expected: Version{
				Major:   0,
				Minor:   5,
				Micro:   1,
				IsPatch: true,
			},
		},
		{
			Version: "v0.5.1-alpha",
			Expected: Version{
				Major:   0,
				Minor:   5,
				Micro:   1,
				IsAlpha: true,
			},
		},
		{
			Version: "v0.5.1-a",
			Expected: Version{
				Major:   0,
				Minor:   5,
				Micro:   1,
				IsAlpha: true,
			},
		},
		{
			Version: "v0.5.1-beta",
			Expected: Version{
				Major:  0,
				Minor:  5,
				Micro:  1,
				IsBeta: true,
			},
		},
		{
			Version: "v0.5.1-b",
			Expected: Version{
				Major:  0,
				Minor:  5,
				Micro:  1,
				IsBeta: true,
			},
		},
		{
			Version: "v0.5.1-RC",
			Expected: Version{
				Major: 0,
				Minor: 5,
				Micro: 1,
				IsRC:  true,
			},
		},

		// Errors

		{
			Version:  "",
			Expected: Version{},
			Error:    fmt.Errorf("version is empty"),
		},
		{
			Version:  "1.0",
			Expected: Version{},
			Error:    fmt.Errorf("version must be in the format [v]X.Y.Z[-suffix]"),
		},
		{
			Version:  "1.0.0-s-s",
			Expected: Version{},
			Error:    fmt.Errorf("version must be in the format [v]X.Y.Z[-suffix]"),
		},
		{
			Version:  "1.0.0-unknown_suffix",
			Expected: Version{},
			Error:    fmt.Errorf("unknown version siffix 'unknown_suffix'"),
		},
		{
			Version:  "1.0.0.0",
			Expected: Version{},
			Error:    fmt.Errorf("version must be in the format [v]X.Y.Z[-suffix]"),
		},
		{
			Version:  "a.1.0",
			Expected: Version{},
			Error:    fmt.Errorf("part 1 ('a') of the version must be a number"),
		},
		{
			Version:  "1.b.0",
			Expected: Version{},
			Error:    fmt.Errorf("part 2 ('b') of the version must be a number"),
		},
		{
			Version:  "1.0.c",
			Expected: Version{},
			Error:    fmt.Errorf("part 3 ('c') of the version must be a number"),
		},
	}

	for _, test := range tests {
		version, err := NewVersion(test.Version)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Errorf("unexpected error")
		}

		if version != nil && *version != test.Expected {
			t.Errorf("version is not correct")
		}
	}
}
