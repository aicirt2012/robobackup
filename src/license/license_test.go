//go:build release

package license

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicense(t *testing.T) {
	text, err := Text()
	assert.Nil(t, err)
	assert.Contains(t, text, "Felix Michel", "license missing")
	assert.Contains(t, text, "github.com/gocarina/gocsv", "transitive license missing")
	assert.True(t, len(text) > 10000, "license text too short")
}
