package fs

import (
	"errors"
	"os"
	"path/filepath"
)

const driveLetters string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// UsedDrives returns an array of all drive letters that are available such as ['C', 'D']!
func UsedDrives() (r []rune) {
	for _, drive := range driveLetters {
		f, err := os.Open(string(drive) + `:\`)
		if err == nil {
			r = append(r, drive)
			f.Close()
		}
	}
	return r
}

// FindUniqueFileAcrossDrives returns drive letter such as 'E' when unique file is found.
func FindUniqueFileAcrossDrives(fileName string, drives []rune) (rune, error) {
	validDrives := findDrivesWithFile(fileName, drives)
	if len(validDrives) == 0 {
		return ' ', errors.New("id file '" + fileName + "' does not exist!")
	} else if len(validDrives) > 1 {
		return ' ', errors.New("id file '" + fileName + "' exist multiple times!")
	}
	return validDrives[0], nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && info.IsDir()
}

func FileNamesWithSuffix(suffix string) ([]string, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pattern := filepath.Join(path, "*"+suffix)
	return filepath.Glob(pattern)
}

func findDrivesWithFile(fileName string, drives []rune) []rune {
	validDrives := []rune{}
	for _, drive := range drives {
		file := string(drive) + `:\\` + fileName
		if !FileExists(file) {
			continue
		}
		validDrives = append(validDrives, drive)
	}
	return validDrives
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return result
}
