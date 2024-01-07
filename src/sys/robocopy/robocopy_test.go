package robocopy

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeFormat = "060102.150405"

func TestRobocopyMirArguments(t *testing.T) {

	cases := []struct {
		name          string
		source        string
		destination   string
		excludeDirs   []string
		dryRun        bool
		forceOverride bool
		expected      string
	}{
		{
			name:          "Simple test",
			source:        `D:\Photos`,
			destination:   `E:\backup.Photos`,
			excludeDirs:   []string{"node_modules"},
			dryRun:        false,
			forceOverride: false,
			expected:      `robocopy D:\Photos E:\backup.Photos /COPY:DAT /DCOPY:T /MIR /B /MT:2 /xD node_modules /R:1000000 /W:30 /NP /UNILOG+:E:\backup.Photos.060102.150405.log /TEE`,
		},
		{
			name:          "Dry run test",
			source:        `D:\Photos`,
			destination:   `E:\backup.Photos`,
			excludeDirs:   []string{"node_modules", ".integrity"},
			dryRun:        true,
			forceOverride: false,
			expected:      `robocopy D:\Photos E:\backup.Photos /COPY:DAT /DCOPY:T /MIR /B /MT:2 /xD node_modules .integrity /R:1000000 /W:30 /NP /UNILOG+:E:\backup.Photos.060102.150405.log /TEE /L`,
		},
		{
			name:          "Simple test with force override",
			source:        `D:\Photos`,
			destination:   `E:\backup.Photos`,
			excludeDirs:   []string{"node_modules"},
			dryRun:        false,
			forceOverride: true,
			expected:      `robocopy D:\Photos E:\backup.Photos /COPY:DAT /DCOPY:T /MIR /B /MT:2 /xD node_modules /IS /IT /R:1000000 /W:30 /NP /UNILOG+:E:\backup.Photos.060102.150405.log /TEE`,
		},
		{
			name:          "Dry run test with force override",
			source:        `D:\Photos`,
			destination:   `E:\backup.Photos`,
			excludeDirs:   []string{"node_modules", ".integrity"},
			dryRun:        true,
			forceOverride: true,
			expected:      `robocopy D:\Photos E:\backup.Photos /COPY:DAT /DCOPY:T /MIR /B /MT:2 /xD node_modules .integrity /IS /IT /R:1000000 /W:30 /NP /UNILOG+:E:\backup.Photos.060102.150405.log /TEE /L`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			startTime := timeNowWithSecondPrecision()
			actual := MIR(c.source, c.destination, c.excludeDirs, c.dryRun, c.forceOverride).String()
			endTime := timeNowWithSecondPrecision()

			assert.Len(t, actual, len(c.expected))
			idx := strings.Index(c.expected, timeFormat)
			actualTimeString := actual[idx:][:len(timeFormat)]
			actualTime, err := time.Parse(timeFormat, actualTimeString)
			assert.NoError(t, err)
			assert.True(t, between(actualTime, startTime, endTime))
			expected := strings.Replace(c.expected, timeFormat, actualTimeString, 1)

			assert.Equal(t, expected, actual)
		})
	}
}

func TestLogFile(t *testing.T) {
	cases := []struct {
		path            string
		expectedLogPath string
	}{
		{
			path:            `E:\backup.Photos`,
			expectedLogPath: `E:\backup.Photos.060102.150405.log`,
		},
		{
			path:            `E:\Backup\Private Desktop\backup.Accounts`,
			expectedLogPath: `E:\Backup\Private Desktop\backup.Accounts.060102.150405.log`,
		},
	}
	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			startTime := timeNowWithSecondPrecision()
			actual := logFile(c.path)
			endTime := timeNowWithSecondPrecision()

			idx := strings.Index(c.expectedLogPath, timeFormat)
			actualTimeString := actual[idx:][:len(timeFormat)]
			actualTime, err := time.Parse(timeFormat, actualTimeString)
			assert.NoError(t, err)
			assert.True(t, between(actualTime, startTime, endTime))

			expected := strings.Replace(c.expectedLogPath, timeFormat, actualTimeString, 1)
			assert.Equal(t, expected, actual)
		})
	}
}

func timeNowWithSecondPrecision() time.Time {
	t, _ := time.Parse(timeFormat, time.Now().Format(timeFormat))
	return t
}

func between(t time.Time, start time.Time, end time.Time) bool {
	return start.Equal(t) || start.Before(t) && t.Before(end) || t.Equal(end)
}
