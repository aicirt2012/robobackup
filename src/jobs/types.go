package jobs

import (
	"strings"

	"github.com/aicirt2012/robobackup/src/config"
)

type Jobs []Job

func (j Jobs) descriptions() (d []string) {
	for _, job := range j {
		d = append(d, job.description())
	}
	return d
}

func NewJobs(dtos config.JobsDto) Jobs {
	jobs := Jobs{}
	for _, dto := range dtos {
		jobs = append(jobs, NewJob(dto))
	}
	return jobs
}

type Job struct {
	Name    string
	Source  Location
	Target  Location
	Options Options
}

func (j *Job) description() string {
	dryRun := ""
	if j.Options.Mir.DryRun {
		dryRun = " [dry run]"
	}
	return j.Source.AbsolutePath() + " => " + j.Target.AbsolutePath() + dryRun
}

func (j *Job) HasResolveErrors() bool {
	return j.Source.HasResolveErrors() || j.Target.HasResolveErrors()
}

func (j Job) ResolveErrorMessage() string {
	msg := []string{j.String()}
	if j.Source.HasResolveErrors() {
		msg = append(msg, "Source "+string(j.Source.ResolveErrors[0])+" could not be resolved")
	}
	if j.Target.HasResolveErrors() {
		msg = append(msg, "Target "+string(j.Target.ResolveErrors[0])+" could not be resolved")
	}
	return strings.Join(msg, "\n")
}

func (j Job) String() string {
	return j.Name + "\n" +
		j.Source.Raw() + " => " + j.Target.Raw()
}

func NewJob(dto config.JobDto) Job {
	return Job{
		Name:    dto.Name,
		Source:  NewLocation(dto.Source),
		Target:  NewLocation(dto.Target),
		Options: NewOptions(dto.Options),
	}
}

type Options struct {
	Mir       MirOptions
	Integrity IntegrityOptions
}

func NewOptions(dto config.OptionsDto) Options {
	return Options{
		Mir:       NewMirOptions(dto.Mir),
		Integrity: NewIntegrityOptions(dto.Integrity),
	}
}

type MirOptions struct {
	DryRun        bool
	ExcludeDirs   []string
	ForceOverride bool
}

func NewMirOptions(dto config.MirOptionsDto) MirOptions {
	return MirOptions{
		DryRun:        dto.DryRun,
		ExcludeDirs:   dto.ExcludeDirs,
		ForceOverride: dto.ForceOverride,
	}
}

type IntegrityOptions struct {
	Upsert bool
	Verify bool
}

func NewIntegrityOptions(dto config.IntegrityOptionsDto) IntegrityOptions {
	return IntegrityOptions{
		Upsert: dto.Upsert,
	}
}

type Location struct {
	raw           string
	DriveId       string
	drive         rune
	Path          string
	ResolveErrors []ResolveError
}

func (l Location) Raw() string {
	return l.raw
}

func NewLocation(raw string) (l Location) {
	l.raw = raw
	sections := strings.Split(raw, "?")
	if len(sections) >= 2 {
		l.DriveId = sections[0]
		l.Path = sections[1]
	}
	return l
}

func (l Location) HasResolveErrors() bool {
	return len(l.ResolveErrors) > 0
}

func (l *Location) AbsolutePath() string {
	return string(l.drive) + ":\\" + strings.TrimLeft(l.Path, "\\")
}

type ResolveError string

const (
	invalidDriveId ResolveError = "driveId"
	invalidPath    ResolveError = "path"
)

type mode string

const (
	Backup mode = "backup"
	Verify mode = "verify"
)
