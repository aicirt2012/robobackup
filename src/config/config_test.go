package config

import (
	"errors"
	"testing"

	"github.com/aicirt2012/robobackup/src/sys/fs"
	"github.com/aicirt2012/robobackup/src/sys/prompt"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestPerformImport(t *testing.T) {
	m1 := mockSelectFile("master.robobackup.yaml", nil)
	defer m1.Unpatch()
	m2 := mockReadFile(ConfigDto{}, nil)
	defer m2.Unpatch()
	config, err := performImport("2.1.0")
	assert.Nil(t, err)
	assert.Equal(t, ConfigDto{}, config)
}

func TestPerformImport_selectFileError(t *testing.T) {
	expectedErr := errors.New("file error")
	m := mockSelectFile("", expectedErr)
	defer m.Unpatch()
	config, err := performImport("2.1.0")
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, ConfigDto{}, config)
}

func TestPerformImport_readFileError(t *testing.T) {
	expectedErr := errors.New("file error")
	m1 := mockSelectFile("master.robobackup.yaml", nil)
	defer m1.Unpatch()
	m2 := mockReadFile(ConfigDto{}, expectedErr)
	defer m2.Unpatch()
	config, err := performImport("2.1")
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, ConfigDto{}, config)
}

func TestSelectFileNames(t *testing.T) {
	cases := []struct {
		name        string
		fileNames   []string
		expected    string
		expectedErr error
	}{
		{
			name:      "One config file exists",
			fileNames: []string{"master.robobackup.yaml"},
			expected:  "master.robobackup.yaml",
		},
		{
			name: "Three config files exists",
			fileNames: []string{
				"m1.robobackup.yaml",
				"m2.robobackup.yaml",
				"m3.robobackup.yaml",
			},
			expected: "m2.robobackup.yaml",
		},
		{
			name:        "No config file exists",
			fileNames:   []string{},
			expectedErr: errors.New("config file '*" + fileSuffix + "' does not exist"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockFsFileNamesWithSuffix(c.fileNames, nil)
			mockPromptSelect(1)
			fileName, err := selectFile()
			assert.Equal(t, c.expectedErr, err)
			assert.Equal(t, c.expected, fileName)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		name           string
		content        string
		version        string
		expectedConfig ConfigDto
		expectedErr    error
	}{
		{
			name:           "Unmarshal error",
			content:        "v: '2.0.0'",
			version:        "2.1.0",
			expectedConfig: ConfigDto{},
			expectedErr:    errors.New("could not deserialize config file: unknown field \"v\""),
		},
		{
			name:    "Version mismatch",
			content: "version: '2.0.0'",
			version: "2.1.0",
			expectedConfig: ConfigDto{
				Version: "2.0.0",
			},
			expectedErr: errors.New("config version does not match the execution version"),
		},
		{
			name:    "Happy case",
			content: "version: '2.1.0'",
			version: "2.1.0",
			expectedConfig: ConfigDto{
				Version: "2.1.0",
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config, err := unmarshal([]byte(c.content), c.version)
			assert.Equal(t, c.expectedErr, err)
			assert.Equal(t, c.expectedConfig, config)
		})
	}
}

func mockSelectFile(result string, err error) *monkey.PatchGuard {
	return monkey.Patch(selectFile, func() (string, error) {
		return result, err
	})
}

func mockReadFile(result ConfigDto, err error) *monkey.PatchGuard {
	return monkey.Patch(readFile, func(string, string) (ConfigDto, error) {
		return result, err
	})
}

func mockFsFileNamesWithSuffix(result []string, err error) *monkey.PatchGuard {
	return monkey.Patch(fs.FileNamesWithSuffix, func(string) ([]string, error) {
		return result, err
	})
}

func mockPromptSelect(result int) *monkey.PatchGuard {
	return monkey.Patch(prompt.Select, func(string, []string) int {
		return result
	})
}
