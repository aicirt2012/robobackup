package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"12345", "54321"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			actual := reverse(c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}
