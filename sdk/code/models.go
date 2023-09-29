package code

import (
	"regexp"
	"strconv"
)

type Analyser struct {
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
		profile: make([]Profile, 0),
	}
}

func (a *Analyser) SetCoverProfile(coverProfile string) *Analyser {
	a.CoverProfile = coverProfile
	a.parseCoverProfile()
	return a
}

func (a *Analyser) SetCoverageByTotal(totalText string) *Analyser {
	reg := regexp.MustCompile(`total:\s+\((\w+)\)\s+(\d+\.\d+)%`)
	matches := reg.FindStringSubmatch(totalText)
	if len(matches) != 3 {
		return a
	}

	//convert to a float
	total, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return a
	}

	a.Coverage = total

	return a
}
