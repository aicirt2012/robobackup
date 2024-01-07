package config

const fileSuffix = ".robobackup.yaml"

type ConfigDto struct {
	Version string  `json:"version"`
	Jobs    JobsDto `json:"jobs"`
}

type JobsDto []JobDto

type JobDto struct {
	Name    string     `json:"name"`
	Source  string     `json:"source"`
	Target  string     `json:"target"`
	Options OptionsDto `json:"options"`
}

type OptionsDto struct {
	Mir       MirOptionsDto       `json:"mir"`
	Integrity IntegrityOptionsDto `json:"integrity"`
}

type MirOptionsDto struct {
	DryRun        bool     `json:"dry-run"`
	ExcludeDirs   []string `json:"exclude-dirs"`
	ForceOverride bool     `json:"force-override"`
}

type IntegrityOptionsDto struct {
	Upsert bool `json:"upsert"`
}
