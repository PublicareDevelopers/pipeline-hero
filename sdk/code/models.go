package code

type Analyser struct {
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
		profile: make([]Profile, 0),
	}
}

func (a *Analyser) SetCoverProfile(coverProfile string) *Analyser {
	a.CoverProfile = coverProfile
	a.parseCoverProfile()
	return a
}
