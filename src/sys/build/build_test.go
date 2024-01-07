//go:build release

package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	actual := Version()
	expected := `^\d+\.\d+\.\d+$`
	assert.Regexp(t, expected, actual, "version pattern invalid")
}
