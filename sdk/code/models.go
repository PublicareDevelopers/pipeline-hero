package code

type Analyser struct {
	GoVersion       string
	Threshold       float64
	Coverage        float64
	Toolchain       string
	CoverProfile    string
	DependencyGraph string
	VulnCheck       string
	Errors          []string
	Warnings        []string
	Profiles        []Profile
	dependencies    []Dependency //Deprecated: we will not use the deep dependency check here anymore; not needed; can be done at platform if needed
}

type Profile struct {
	Folder   string
	Coverage float64
	Duration float64
}

type Dependency struct {
	From      string
	To        string
	Version   string
	Updatable bool
	UpdateTo  string
}

func NewAnalyser() *Analyser {
	return &Analyser{
		Threshold: 75.0,
		Profiles:  make([]Profile, 0),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
	}
}

func (a *Analyser) SetThreshold(threshold float64) *Analyser {
	a.Threshold = threshold
	return a
}

func (a *Analyser) SetCoverProfile(coverProfile string) *Analyser {
	a.CoverProfile = coverProfile
	a.parseCoverProfile()
	return a
}

func (a *Analyser) SetDependencyGraph(dependencyGraph string) *Analyser {
	a.DependencyGraph = dependencyGraph
	a.parseDependencyGraph()
	return a
}

func (a *Analyser) SetVulnCheck(vulnCheck string) *Analyser {
	a.VulnCheck = vulnCheck
	return a
}

func (a *Analyser) SetGoVersion(goVersion string) *Analyser {
	a.GoVersion = goVersion
	return a
}

func (a *Analyser) PushError(err string) *Analyser {
	a.Errors = append(a.Errors, err)
	return a
}

func (a *Analyser) PushWarning(warning string) *Analyser {
	a.Warnings = append(a.Warnings, warning)
	return a
}

func (a *Analyser) GetErrors() []string {
	return a.Errors
}

func (a *Analyser) GetWarnings() []string {
	return a.Warnings
}
