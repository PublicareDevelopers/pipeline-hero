package code

import "sync"

type Analyser struct {
	GoVersion        string
	HasGoVersionFail bool
	Threshold        float64
	Coverage         float64
	HasCoverageFail  bool
	Toolchain        string
	CoverProfile     string
	DependencyGraph  string
	VulnCheck        string
	HasVuln          bool
	HasErrors        bool
	errors           []string
	warnings         []string
	profiles         []Profile
	dependencies     []Dependency
	lock             *sync.Mutex
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
		profiles:  make([]Profile, 0),
		errors:    make([]string, 0),
		warnings:  make([]string, 0),
		lock:      &sync.Mutex{},
	}
}

func (a *Analyser) SetThreshold(threshold float64) *Analyser {
	a.lock.Lock()
	a.Threshold = threshold
	a.lock.Unlock()
	return a
}

func (a *Analyser) SetCoverProfile(coverProfile string) *Analyser {
	a.lock.Lock()
	a.CoverProfile = coverProfile
	a.lock.Unlock()
	a.parseCoverProfile()
	return a
}

func (a *Analyser) SetDependencyGraph(dependencyGraph string) *Analyser {
	a.lock.Lock()
	a.DependencyGraph = dependencyGraph
	a.lock.Unlock()
	a.parseDependencyGraph()
	return a
}

func (a *Analyser) SetVulnCheck(vulnCheck string) *Analyser {
	a.lock.Lock()
	a.VulnCheck = vulnCheck
	a.lock.Unlock()
	return a
}

func (a *Analyser) SetGoVersion(goVersion string) *Analyser {
	a.lock.Lock()
	a.GoVersion = goVersion
	a.lock.Unlock()
	return a
}

func (a *Analyser) PushError(err string) *Analyser {
	a.lock.Lock()
	a.errors = append(a.errors, err)
	a.lock.Unlock()
	return a
}

func (a *Analyser) PushWarning(warning string) *Analyser {
	a.lock.Lock()
	a.warnings = append(a.warnings, warning)
	a.lock.Unlock()
	return a
}

func (a *Analyser) GetErrors() []string {
	return a.errors
}

func (a *Analyser) GetWarnings() []string {
	return a.warnings
}
