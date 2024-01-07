package robocopy

import (
	"path/filepath"
	"time"
)

// ------------------------------------------------------------------------------
// ---------Example Robocopy output configuration--------------------------------
// ------------------------------------------------------------------------------
//
//   ROBOCOPY     ::     Robust File Copy for Windows
//-------------------------------------------------------------------------------
//
//  Started : Wednesday, 19 February 2020 02:57:37
//   Source : D:\Photo\
//     Dest : Y:\backup.Photo\
//
//    Files : *.*
//
//  Options : *.* /TEE /S /E /DCOPY:T /COPY:DAT /PURGE /MIR /B /NP /MT:2 /R:1000000 /W:30

// MIR initializes a new robocopy MIR operation
func MIR(source string, destination string, excludeDirs []string, dryRun bool, forceOverride bool) *rc {
	rc := rc{
		source:      source,
		destination: destination,
		files:       []string{"*.*"},
		copy: copy{
			mirror:        true,
			fileInfo:      "DAT",
			dirInfo:       "T",
			backupMode:    true,
			multiThreaded: 2,
		},
		file: file{
			copyArchived:    false,
			excludeDirs:     excludeDirs,
			includeSameSize: forceOverride,
			includeTweaked:  forceOverride,
		},
		retry: retry{
			count: 1000000,
			wait:  30,
		},
		logging: logging{
			hideProgress: true,
			showEta:      false,
			output: output{
				unicode:   true,
				overwrite: true,
				file:      logFile(destination),
			},
			showAndLog: true,
		},
		dryRun: dryRun,
	}
	return &rc
}

// Uses the destination parent folder as location for the log file!
func logFile(path string) string {
	dateTime := time.Now().Format("060102.150405")
	parentPath := filepath.Dir(path)
	dir := filepath.Base(path)
	filename := dir + "." + dateTime + ".log"
	return filepath.Join(parentPath, filename)
}
