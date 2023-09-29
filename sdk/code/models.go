package code

type Analyser struct {
	Threshold    float64
	Coverage     float64
	CoverProfile string
	profile      []Profile
}

type Profile struct {
	Folder   string
	Coverage float64
	Duration float64
}

func NewAnalyser() *Analyser {
	return &Analyser{
		Threshold: 75.0,
		profile:   make([]Profile, 0),
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
