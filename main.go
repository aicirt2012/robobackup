package main

import (
	"log"
	"os"

	"github.com/aicirt2012/robobackup/src/config"
	"github.com/aicirt2012/robobackup/src/jobs"
	"github.com/aicirt2012/robobackup/src/license"
	"github.com/aicirt2012/robobackup/src/sys"
	"github.com/aicirt2012/robobackup/src/sys/build"
)

func main() {
	log.SetFlags(log.Lmsgprefix)

	sys.PrintHeadline("RoboBackup Utility", build.Info())
	license.Print(os.Args)

	sys.AssertAdminPermissions()
	cfg := config.Import(build.Version())
	jobs.Process(cfg.Jobs)
}
