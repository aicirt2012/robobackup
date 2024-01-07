package build

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"time"
)

var version = "develop"

func Version() string {
	return version
}

func Info() string {
	b := build()
	return fmt.Sprintf("v%v %v %v", version, b.commitHash, b.fmtCommitTime())
}

type info struct {
	commitHash string
	commitTime time.Time
	modified   bool
}

func (i info) fmtCommitTime() string {
	return i.commitTime.Local().Format("2006-01-02 15:04:05")
}

func build() (i info) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				i.commitHash = setting.Value
			case "vcs.time":
				t, err := time.Parse(time.RFC3339, setting.Value)
				if err != nil {
					log.Fatal("could not parse vcs.time")
				}
				i.commitTime = t
			case "vcs.modified":
				m, err := strconv.ParseBool(setting.Value)
				if err != nil {
					log.Fatal("could not parse vcs.modified")
				}
				i.modified = m
			}
		}
	}
	return i
}
