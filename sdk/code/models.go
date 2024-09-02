package code

import (
	"encoding/json"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/npm"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/php"
	"sync"
)

type Analyser struct {
	Module           GoMod
	GoVersion        string //will have the OS also here
	HasGoVersionFail bool
	Threshold        float64
	Coverage         float64
	HasCoverageFail  bool
	CoverProfile     string
	DependencyGraph  string
	Updates          []RequireUpdate
	VulnCheck        string
	HasVuln          bool
	HasErrors        bool
	Errors           []string
	Warnings         []string
	TestResult       string
	Profiles         []Profile
	dependencies     []Dependency //Deprecated: use Updates instead
	lock             *sync.Mutex
}

type JSAnalyser struct {
	Threshold        float64
	Coverage         float64
	HasCoverageFail  bool
	VulnCheck        string
	HasVulnCheckFail bool
	HasErrors        bool
	Errors           []string
	Warnings         []string
	HasWarnings      bool
	TestResult       string
	OutDates         []npm.OutDate
	lock             *sync.Mutex
}

type PHPAnalyser struct {
	Threshold        float64
	Coverage         float64
	HasCoverageFail  bool
	VulnCheck        string
	HasVulnCheckFail bool
	HasErrors        bool
	Errors           []string
	Warnings         []string
	HasWarnings      bool
	TestResult       string
	OutDates         []php.OutDate
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

type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module    ModPath
	Go        string
	Toolchain string
	Require   []Require
	Exclude   []Module
	Replace   []Replace
	Retract   []Retract
}

type ModPath struct {
	Path       string
	Deprecated string
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type RequireUpdate struct {
	Path             string
	Version          string
	AvailableVersion string
	Indirect         bool
}

type Replace struct {
	Old Module
	New Module
}

type Retract struct {
	Low       string
	High      string
	Rationale string
}

func NewAnalyser() *Analyser {
	return &Analyser{
		Threshold: 75.0,
		Profiles:  make([]Profile, 0),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
		Updates:   make([]RequireUpdate, 0),
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

func (a *Analyser) SetModule() *Analyser {
	moduleJson, err := cmds.AnalyseModule()
	if err != nil {
		a.lock.Lock()
		a.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot find the module: %s\n", err))
		a.lock.Unlock()
		return a
	}
	module := GoMod{}
	err = json.Unmarshal([]byte(moduleJson), &module)
	if err != nil {
		a.lock.Lock()
		a.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot parse the module: %s\n", err))
		a.lock.Unlock()
		return a
	}

	a.lock.Lock()
	a.Module = module
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
	a.Errors = append(a.Errors, err)
	a.lock.Unlock()
	return a
}

func (a *Analyser) PushWarning(warning string) *Analyser {
	a.lock.Lock()
	a.Warnings = append(a.Warnings, warning)
	a.lock.Unlock()
	return a
}

func (a *Analyser) GetErrors() []string {
	return a.Errors
}

func (a *Analyser) GetWarnings() []string {
	return a.Warnings
}

func (a *Analyser) GetProfiles() []Profile {
	return a.Profiles
}
