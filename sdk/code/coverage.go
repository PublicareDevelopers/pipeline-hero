package code

import (
	"fmt"
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

func (a *Analyser) GetCoverageInterpretation() string {
	if a.Coverage < BAD {
		return fmt.Sprintf("coverage is BAD, have %.2f  percent\n", a.Coverage)
	}

	if a.Coverage < MEDIUM && a.Coverage >= BAD {
		return fmt.Sprintf("coverage is ok, have %.2f  percent\n", a.Coverage)
	}

	if a.Coverage >= MEDIUM {
		return fmt.Sprintf("coverage is good, have %.2f  percent\n", a.Coverage)
	}

	return fmt.Sprintf("coverage is unknown\n")
}
