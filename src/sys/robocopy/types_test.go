package robocopy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertExecutable(t *testing.T) {

	cases := []struct {
		name        string
		source      string
		destination string
		expected    error
	}{
		{
			name:        "Too short source",
			source:      `E:\`,
			destination: `F:\MyPhotos`,
			expected:    errors.New("[ROBOCOPY] Source or destination path to short"),
		},
		{
			name:        "Too short destination",
			source:      `E:\MyPhotos`,
			destination: `F:\`,
			expected:    errors.New("[ROBOCOPY] Source or destination path to short"),
		},
		{
			name:        "Similar source and destination",
			source:      `E:\MyPhotos`,
			destination: `E:\MyPhotos`,
			expected:    errors.New("[ROBOCOPY] Source and destination must not be similar"),
		},
		{
			name:        "Source not exists",
			source:      `JJ:\MyPhotos`,
			destination: `F:\MyPhotos`,
			expected:    errors.New("[ROBOCOPY] Source does not exist"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rc := rc{
				source:      c.source,
				destination: c.destination,
			}
			actual := rc.assertExecutable()
			assert.Equal(t, c.expected, actual)
		})
	}
}
