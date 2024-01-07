package tests

import (
	"testing"

	"github.com/aicirt2012/robobackup/src/config"
	"github.com/aicirt2012/robobackup/tests/assert"
	"github.com/aicirt2012/robobackup/tests/execute"
	"github.com/aicirt2012/robobackup/tests/setup"
)

var backupInput = []string{"0", "0", "186", "exits"}
var verifyInput = []string{"1", "0", "1", "186", "exit"}
var plainOptions = config.OptionsDto{}
var dryRunOptions = config.OptionsDto{
	Mir: config.MirOptionsDto{
		DryRun: true,
	},
}
var excludeDirsOptions = config.OptionsDto{
	Mir: config.MirOptionsDto{
		ExcludeDirs: []string{
			".git",
			"node_modules",
		},
	},
}
var upsertOptions = config.OptionsDto{
	Integrity: config.IntegrityOptionsDto{
		Upsert: true,
	},
}

func TestNewFile(t *testing.T) {
	t.Parallel()
	newFile := setup.NewFile(`a.txt`, `2023-05-06T00:40:21+02:00`, `a sample txt`)
	s := setup.NewScenario("mir.newFile",
		setup.Files{newFile},
		setup.Files{},
		plainOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{newFile},
		s.TargetDir, setup.Files{newFile})
}

func TestNewerFile(t *testing.T) {
	t.Parallel()
	newFile := setup.NewFile(`a.txt`, `2023-05-06T00:40:21+02:00`, `a sample txt`)
	oldFile := setup.NewFile(`a.txt`, `2022-05-06T00:40:21+02:00`, `a sample txt`)
	s := setup.NewScenario("mir.newerFile",
		setup.Files{newFile},
		setup.Files{oldFile},
		plainOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{newFile},
		s.TargetDir, setup.Files{newFile})
}

func TestOlderFile(t *testing.T) {
	t.Parallel()
	newFile := setup.NewFile(`a.txt`, `2023-05-06T00:40:21+02:00`, `a sample txt`)
	oldFile := setup.NewFile(`a.txt`, `2022-05-06T00:40:21+02:00`, `a sample txt`)
	s := setup.NewScenario("mir.olderFile",
		setup.Files{oldFile},
		setup.Files{newFile},
		plainOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{oldFile},
		s.TargetDir, setup.Files{oldFile})
}

func TestDeletedFile(t *testing.T) {
	t.Parallel()
	deletedFile := setup.NewFile(`a.txt`, `2022-05-06T00:40:21+02:00`, `a sample txt`)
	s := setup.NewScenario("mir.deletedFile",
		setup.Files{},
		setup.Files{deletedFile},
		plainOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{},
		s.TargetDir, setup.Files{})
}

func TestNesting(t *testing.T) {
	t.Parallel()
	existingFile := setup.NewFile(`dir\a.txt`, `2022-05-06T00:40:21+02:00`, `a sample txt`)
	newFile := setup.NewFile(`dir\b.txt`, `2022-05-06T00:40:21+02:00`, `b sample txt`)
	s := setup.NewScenario("mir.nesting",
		setup.Files{
			existingFile,
			newFile,
		},
		setup.Files{
			existingFile,
		},
		plainOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{
		existingFile,
		newFile,
	}, s.TargetDir, setup.Files{
		existingFile,
		newFile,
	})
}

func TestDryRun(t *testing.T) {
	t.Parallel()
	newFile := setup.NewFile(`a.txt`, `2022-05-06T00:40:21+02:00`, `a sample txt`)
	s := setup.NewScenario("mir.dryRun",
		setup.Files{newFile},
		setup.Files{},
		dryRunOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{}, s.TargetDir, setup.Files{})
}

func TestExcludeDirs(t *testing.T) {
	t.Parallel()
	newFile := setup.NewFile(`a.txt`, `2022-05-06T00:40:21+02:00`, `a sample txt`)
	gitFile := setup.NewFile(`x\.git\a.txt`, `2022-05-06T00:40:21+02:00`, `any content`)
	nodeFile := setup.NewFile(`y\node_modules\a.txt`, `2022-05-06T00:40:21+02:00`, `any content`)
	s := setup.NewScenario("mir.excludeDirs",
		setup.Files{newFile, gitFile, nodeFile},
		setup.Files{},
		excludeDirsOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)

	assert.FileSystem(t, s.TargetDir, setup.Files{newFile}, s.TargetDir, setup.Files{newFile})
}

func TestBackupAndVerify(t *testing.T) {
	t.Parallel()
	newFile := setup.NewFile(`a.txt`, `2023-05-06T00:40:21+02:00`, `a sample txt`)
	s := setup.NewScenario("backupAndVerify",
		setup.Files{newFile},
		setup.Files{},
		upsertOptions,
	)

	execute.Binary(s.BaseDir, s.Binary, backupInput...)
	execute.Binary(s.BaseDir, s.Binary, verifyInput...)
	assert.FileSystem(t, s.SourceDir, setup.Files{newFile}, s.TargetDir, setup.Files{newFile})
	assert.LastVerifyLog(t, s.TargetDir, 1, 0)
}
