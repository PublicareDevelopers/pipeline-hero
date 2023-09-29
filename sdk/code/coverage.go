package code

import (
	"regexp"
	"strconv"
)

const BAD = 50.0
const MEDIUM = 75.0

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
