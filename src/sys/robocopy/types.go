package robocopy

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aicirt2012/robobackup/src/sys/fs"
)

const robocopy = "robocopy"
const logPrefix = "[ROBOCOPY] "
const backupFolderPrefix = "backup."

type rc struct {
	source      string
	destination string
	files       []string
	copy        copy
	file        file
	retry       retry
	logging     logging
	dryRun      bool
}

type copy struct {
	fileInfo      string // Specifies the file properties to be copied. [/COPY:<CopyFlags>] e.g. 'DAT'
	dirInfo       string // Specifies the file properties to be copied. [/DCOPY:<CopyFlags>] e.g. 'DAT'
	mirror        bool   // /MIR deletes on destination if deleted on source
	backupMode    bool   // Copies files in Backup mode. [/B]
	multiThreaded int    // Creates multi-threaded copies with N threads. N must be an integer between
}

type file struct {
	copyArchived    bool
	excludeDirs     []string
	includeSameSize bool // Includes the same files. Same files are identical in name, size, times, and all attributes.
	includeTweaked  bool // Includes "tweaked" files. Tweaked files have the same name, size, and times, but different attributes.
}

type retry struct {
	count int // /R
	wait  int // /W
}

type logging struct {
	hideProgress bool // /NP
	showEta      bool // /ETA
	output       output
	showAndLog   bool // /TEE
}

type output struct {
	file      string
	overwrite bool
	unicode   bool
}

// Execute executes the declared robocopy MIR operation!
func (rc *rc) Execute() error {
	if err := rc.assertExecutable(); err != nil {
		return err
	}

	cmd := exec.Command(robocopy, rc.args()...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			log.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}
		log.Println(logPrefix + "Successfully completed!")
	}
	return nil
}

func (rc *rc) args() []string {
	args := []string{}

	args = append(args, rc.source)
	args = append(args, rc.destination)

	// Copy options
	args = append(args, "/COPY:"+rc.copy.fileInfo)
	args = append(args, "/DCOPY:"+rc.copy.dirInfo)
	args = condAppend(rc.copy.mirror, args, "/MIR")
	args = condAppend(rc.copy.backupMode, args, "/B")
	args = append(args, "/MT:"+strconv.Itoa(rc.copy.multiThreaded))

	// File options
	args = condAppend(rc.file.copyArchived, args, "/A")
	if len(rc.file.excludeDirs) > 0 {
		args = append(args, "/xD")
		args = append(args, rc.file.excludeDirs...)
	}
	args = condAppend(rc.file.includeSameSize, args, "/IS")
	args = condAppend(rc.file.includeTweaked, args, "/IT")

	// Retry options
	args = append(args, "/R:"+strconv.Itoa(rc.retry.count))
	args = append(args, "/W:"+strconv.Itoa(rc.retry.wait))

	// Logging options
	args = condAppend(rc.logging.hideProgress, args, "/NP")
	args = condAppend(rc.logging.showEta, args, "/ETA")
	if (rc.logging.output != output{}) {
		arg := "/"
		if rc.logging.output.unicode {
			arg += "UNI"
		}
		arg += "LOG"
		if rc.logging.output.overwrite {
			arg += "+"
		}
		arg += ":" + rc.logging.output.file
		args = append(args, arg)
	}
	args = condAppend(rc.logging.showAndLog, args, "/TEE")
	args = condAppend(rc.dryRun, args, "/L")
	fmt.Printf("%v", args)
	return args
}

func condAppend(condition bool, args []string, value string) []string {
	if condition {
		args = append(args, value)
	}
	return args
}

func (rc *rc) String() string {
	return robocopy + " " + strings.Join(rc.args(), " ")
}

func (rc *rc) assertExecutable() error {
	if len(rc.source) < 4 || len(rc.destination) < 4 {
		return errors.New(logPrefix + "Source or destination path to short")
	}
	if rc.source == rc.destination {
		return errors.New(logPrefix + "Source and destination must not be similar")
	}
	if !fs.DirExists(rc.source) {
		return errors.New(logPrefix + "Source does not exist")
	}
	if !fs.DirExists(rc.destination) {
		return errors.New(logPrefix + "Destination does not exist")
	}
	if dir := filepath.Base(rc.destination); !strings.HasPrefix(dir, backupFolderPrefix) {
		return errors.New(logPrefix + "Destination does not have expected backup folder prefix '" + backupFolderPrefix + "'")
	}
	return nil
}
