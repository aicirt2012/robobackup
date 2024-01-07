package jobs

import (
	"log"
	"strings"

	"github.com/aicirt2012/robobackup/src/config"
	"github.com/aicirt2012/robobackup/src/sys"
	"github.com/aicirt2012/robobackup/src/sys/fs"
	"github.com/aicirt2012/robobackup/src/sys/prompt"
	"github.com/aicirt2012/robobackup/src/sys/robocopy"

	"github.com/aicirt2012/fileintegrity"
)

func Process(dtos config.JobsDto) {
	jobs := NewJobs(dtos)
	resolvedJobs := ResolveJobs(jobs)
	switch SelectMode() {
	case Backup:
		ExecuteBackup(SelectBackup(resolvedJobs))
	case Verify:
		VerifyPaths(SelectVerifyPaths(resolvedJobs))
	}
}

func ResolveJobs(jobs Jobs) (resolved Jobs) {
	driveLetters := fs.UsedDrives()
	for _, job := range jobs {
		ResolveLocation(&job.Source, driveLetters)
		ResolveLocation(&job.Target, driveLetters)
		if job.HasResolveErrors() {
			log.Println("\n" + job.ResolveErrorMessage())
			continue
		}
		resolved = append(resolved, job)
	}
	if len(resolved) == 0 {
		log.Println("\nNo executable job found!")
		prompt.PressAnyKeyToClose()
		log.Fatal()
	}
	return resolved
}

func ResolveLocation(l *Location, driveLetters []rune) {
	drive, sourceErr := fs.FindUniqueFileAcrossDrives(l.DriveId, driveLetters)
	if sourceErr != nil {
		l.ResolveErrors = append(l.ResolveErrors, invalidDriveId)
		return
	}
	l.drive = drive
	if !fs.DirExists(l.AbsolutePath()) {
		l.ResolveErrors = append(l.ResolveErrors, invalidPath)
	}
}

func SelectMode() mode {
	options := []string{"Backup", "Manual verification"}
	i := prompt.Select("Please select mode:", options)
	switch i {
	case 0:
		return Backup
	case 1:
		return Verify
	}
	return Verify
}

func SelectBackup(jobs Jobs) (selected Jobs) {
	options := append(jobs.descriptions(), "All")
	i := prompt.Select("Please select jobs to be executed:", options)

	// Select all jobs or a dedicated job
	if i == len(jobs) {
		selected = append(selected, jobs...)
	} else {
		selected = append(selected, jobs[i])
	}

	log.Println("\nFollowing jobs are selected for execution:")
	for _, job := range selected {
		log.Println(job.description())
	}
	log.Println()

	prompt.ConfirmWithNumber()
	return selected
}

func ExecuteBackup(jobs Jobs) {
	for _, job := range jobs {
		sys.PrintHeadline(job.Name, job.description())
		if job.Options.Integrity.Upsert {
			err := fileintegrity.Upsert(job.Source.AbsolutePath(), logOptions(false))
			if err != nil {
				log.Println(err)
			}
		}
		err := robocopy.MIR(
			job.Source.AbsolutePath(),
			job.Target.AbsolutePath(),
			job.Options.Mir.ExcludeDirs,
			job.Options.Mir.DryRun,
			job.Options.Mir.ForceOverride,
		).Execute()
		if err != nil {
			log.Println(err)
		}
	}
	prompt.PressAnyKeyToClose()
}

func SelectVerifyPaths(jobs Jobs) (paths []string) {
	jobOptions := append(jobs.descriptions(), "All")
	i := prompt.Select("Please select jobs to be verified:", jobOptions)

	// Select all jobs or a dedicated job
	var pathOptions [][]string
	if i == len(jobs) {
		sourceOptions := []string{}
		targetOptions := []string{}
		for _, job := range jobs {
			sourceOptions = append(sourceOptions, job.Source.AbsolutePath())
			targetOptions = append(targetOptions, job.Target.AbsolutePath())
		}
		pathOptions = [][]string{
			sourceOptions,
			targetOptions,
		}
	} else {
		pathOptions = [][]string{
			{jobs[i].Source.AbsolutePath()},
			{jobs[i].Target.AbsolutePath()},
		}
	}

	i = prompt.SelectMultiLineOptions("Please select source or target paths to be verified:", pathOptions)
	paths = pathOptions[i]

	log.Println("\nFollowing job paths are selected for verification:")
	log.Println(strings.Join(pathOptions[i], "\n"))
	log.Println()

	prompt.ConfirmWithNumber()
	return paths
}

func VerifyPaths(paths []string) {
	for _, path := range paths {
		fileintegrity.Verify(path, logOptions(true))
	}
	prompt.PressAnyKeyToClose()
}

func logOptions(file bool) fileintegrity.Options {
	return fileintegrity.Options{
		LogConsole:  true,
		LogFile:     file,
		Backup:      false,
		ProgressBar: false,
	}
}
