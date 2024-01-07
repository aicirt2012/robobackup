package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		expected Location
	}{
		{
			name:     "Happy case",
			raw:      ".driveId?/path/to/location",
			expected: Location{raw: ".driveId?/path/to/location", DriveId: ".driveId", Path: "/path/to/location"},
		},
		{
			name:     "Empty case",
			raw:      "",
			expected: Location{raw: "", DriveId: "", Path: ""},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := NewLocation(c.raw)
			assert.Equal(t, c.expected.raw, actual.raw)
			assert.Equal(t, c.expected.DriveId, actual.DriveId)
			assert.Equal(t, c.expected.drive, actual.drive)
			assert.Equal(t, c.expected.Path, actual.Path)
			assert.Equal(t, c.expected.ResolveErrors, actual.ResolveErrors)
		})
	}
}
