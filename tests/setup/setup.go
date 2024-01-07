package setup

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aicirt2012/robobackup/src/config"
	"github.com/aicirt2012/robobackup/src/sys"
	"github.com/aicirt2012/robobackup/src/sys/build"

	"sigs.k8s.io/yaml"
)

const testDriveId = `.testdrive`
const sourceDir = `source`
const backupDir = `backup.source`

func NewScenario(name string, source Files, target Files, options config.OptionsDto) Scenario {
	sys.AssertAdminPermissions()
	s := Scenario{}
	s.BaseDir = upsertScenarioDir(name)
	s.Binary = buildExecutable(s.BaseDir)
	createConfigFile(s.BaseDir, options)
	upsertDriveId(s.BaseDir)
	s.SourceDir = createTestDir(s.BaseDir, sourceDir, source)
	s.TargetDir = createTestDir(s.BaseDir, backupDir, target)
	return s
}

func upsertScenarioDir(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not get working dir", err)
	}
	dir := filepath.Join(wd, "cases", name)
	err = os.RemoveAll(dir)
	if err != nil {
		log.Fatal("could not wipe", err)
	}
	err = os.MkdirAll(dir, 0644)
	if err != nil {
		log.Fatal("could not create scenario dir", err)
	}
	return dir
}

func createTestDir(baseDir string, dir string, files Files) string {
	absoluteDir := filepath.Join(baseDir, dir)
	if err := os.Mkdir(absoluteDir, 0644); err != nil {
		log.Fatal("could create test dir", err)
	}
	err := createTestFiles(absoluteDir, files)
	if err != nil {
		log.Fatal("could not create test files", err)
	}
	return absoluteDir
}

func createTestFiles(baseDir string, files []File) error {
	for _, file := range files {
		absoluteFilename := filepath.Join(baseDir, file.RelativePath)
		absoluteDir := filepath.Dir(absoluteFilename)
		if err := os.MkdirAll(absoluteDir, 0644); err != nil {
			return err
		}
		if err := os.WriteFile(absoluteFilename, []byte(file.Content), file.Permission); err != nil {
			return err
		}
		if err := os.Chtimes(absoluteFilename, file.ModTime, file.ModTime); err != nil {
			return err
		}
	}
	return nil
}

func buildExecutable(dir string) string {
	filename := filepath.Join(dir, "robobackup.exe")
	cmd := exec.Command("go", "build", "-tags", "test", "-o", filename, "../main.go")
	err := cmd.Run()
	if err != nil {
		log.Fatal("error building test executable", err.Error())
	}
	return filename
}

func createConfigFile(basePath string, options config.OptionsDto) string {
	relBasePath := extractVolume(basePath)
	filename := filepath.Join(basePath, "test.robobackup.yaml")
	config := config.ConfigDto{
		Version: build.Version(),
		Jobs: []config.JobDto{
			{
				Name:    `Test`,
				Source:  filepath.Join(testDriveId+`?`, relBasePath, sourceDir),
				Target:  filepath.Join(testDriveId+`?`, relBasePath, backupDir),
				Options: options,
			},
		},
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		log.Fatal("could create config file", err)
	}
	return string(content)
}

func upsertDriveId(path string) {
	volume := filepath.VolumeName(path)
	filename := volume + string(filepath.Separator) + testDriveId
	err := touchFile(filename)
	if os.IsExist(err) {
		return
	}
	if err != nil {
		log.Fatal("could not create test drive id file", err)
	}
}

// Normalizes static test data to fit the os specific notation
func NormalizePath(path string) string {
	path = strings.ReplaceAll(path, `\`, string(filepath.Separator))
	return strings.ReplaceAll(path, `/`, string(filepath.Separator))
}

func extractVolume(absolutePath string) string {
	volume := filepath.VolumeName(absolutePath)
	return strings.TrimPrefix(absolutePath, volume)
}

func touchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func parseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Fatal("could not parse test data modTime ", value)
	}
	return t
}
