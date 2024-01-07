package setup

import (
	"io/fs"
	"time"
)

const Integrity string = ".integrity"

type Scenario struct {
	BaseDir     string
	SourceDir   string
	SourceFiles Files
	TargetDir   string
	TargetFiles Files
	Binary      string
}

type File struct {
	RelativePath string
	ModTime      time.Time
	Permission   fs.FileMode
	Content      string
}

func permission(readOnly bool) fs.FileMode {
	if readOnly {
		return 0444
	}
	return 0644
}

func NewFile(relativePath string, modTime string, content string) File {
	return File{
		RelativePath: NormalizePath(relativePath),
		ModTime:      parseTime(modTime),
		Permission:   permission(false),
		Content:      content,
	}
}

type Files []File

func (f *Files) Map() map[string]File {
	m := make(map[string]File)
	for _, i := range *f {
		m[i.RelativePath] = i
	}
	return m
}
