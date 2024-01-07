package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVolume(t *testing.T) {
	expected := `\dir\sub`
	actual := extractVolume(`D:\dir\sub`)
	assert.Equal(t, expected, actual)
}
