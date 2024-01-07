package assert

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"

	"github.com/aicirt2012/robobackup/tests/setup"

	"github.com/stretchr/testify/assert"
)

func FileSystem(t *testing.T, sourceDir string, expectedSourceFiles setup.Files, targetDir string, expectedTargetFiles setup.Files) {
	FilesExist(t, sourceDir, expectedSourceFiles)
	FilesExist(t, targetDir, expectedTargetFiles)
}

func FilesExist(t *testing.T, basePath string, expectedFiles setup.Files) {
	actualFiles, err := listFiles(basePath)
	if err != nil {
		assert.Fail(t, "could not list actual files", err)
	}
	for _, ef := range expectedFiles {
		af := actualFiles.Map()[ef.RelativePath]
		assert.Equal(t, ef.RelativePath, af.RelativePath, "relative path")
		assert.Equal(t, ef.ModTime.UTC(), af.ModTime.UTC(), "mode time")
		assert.Equal(t, ef.Content, af.Content, "content")
	}
	assert.Equal(t, len(expectedFiles), len(actualFiles))
}

func LastVerifyLog(t *testing.T, basePath string, valid int, invalid int) {
	pattern := filepath.Join(basePath, setup.Integrity, "*.verify.log")
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal("could list log files", err)
	}
	if len(files) == 0 {
		log.Fatal("no matching log file found")
	}
	slices.Sort(files)
	b, err := os.ReadFile(files[len(files)-1])
	if err != nil {
		log.Fatal("could not read log file", err)
	}
	content := string(b)
	assert.Contains(t, content, "Verified valid files:                    "+strconv.Itoa(valid), "valid files")
	assert.Contains(t, content, "Verified invalid files:                  "+strconv.Itoa(invalid), "invalid files")
}

func listFiles(basePath string) (files setup.Files, err error) {
	e := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if basePath == path {
			return nil
		}
		if info.IsDir() && info.Name() == setup.Integrity {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relativeFilename, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		modTime := info.ModTime()
		files = append(files, setup.File{
			RelativePath: setup.NormalizePath(relativeFilename),
			ModTime:      modTime,
			Permission:   info.Mode().Perm(),
			Content:      string(content),
		})
		return nil
	})
	return files, e
}
