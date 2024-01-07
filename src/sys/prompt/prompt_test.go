package prompt

import (
	"os"
	"strings"
	"testing"

	"github.com/aicirt2012/robobackup/src/sys/prompt/rand"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestSelectMultiLineOptions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			input:    "0\n",
			expected: 0,
		},
		{
			input:    "1\n",
			expected: 1,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reader = strings.NewReader(c.input) // Mock reader
			options := [][]string{
				{"Option 1", "Label 1"},
				{"Option 2", "Label 2"},
			}
			actual := SelectMultiLineOptions("Please select?", options)
			assert.Equal(t, c.expected, actual, c.name)
		})
	}
}

func TestSelect(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			input:    "0\n",
			expected: 0,
		},
		{
			input:    "1\n",
			expected: 1,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reader = strings.NewReader(c.input) // Mock reader
			options := []string{"Option 1", "Option 2"}
			actual := Select("Please select?", options)
			assert.Equal(t, c.expected, actual, c.name)
		})
	}
}

func TestConfirmWithYes(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			input:    "y\n",
			expected: true,
		},
		{
			input:    "Y\n",
			expected: true,
		},
		{
			input:    "n\n",
			expected: false,
		},
		{
			input:    "N\n",
			expected: false,
		},
		{
			input:    "x\n",
			expected: false,
		},
		{
			input:    "\n",
			expected: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reader = strings.NewReader(c.input) // Mock reader
			actual := ConfirmWithYes("Any question?")
			assert.Equal(t, c.expected, actual, c.name)
		})
	}
}

func TestConfirmWithNumber_success(t *testing.T) {
	reader = strings.NewReader("12\n") // Mock reader
	patch := monkey.Patch(rand.ThreeDigits, func() string {
		return "12"
	})
	defer patch.Unpatch()
	assert.NotPanics(t, ConfirmWithNumber)
}

func TestConfirmWithNumber_exit(t *testing.T) {
	reader = strings.NewReader("no-number\n") // Mock reader
	patch := monkey.Patch(os.Exit, func(int) {
		panic("os.Exit")
	})
	defer patch.Unpatch()
	assert.PanicsWithValue(t, "os.Exit", ConfirmWithNumber)
}
