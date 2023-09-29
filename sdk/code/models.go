package code

type Analyser struct {
	Threshold       float64
	Coverage        float64
	CoverProfile    string
	DependencyGraph string
	profiles        []Profile
	dependencies    []Dependency
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
