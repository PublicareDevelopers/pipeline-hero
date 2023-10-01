package code

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (a *Analyser) parseCoverProfile() {
	//read a.CoverProfile line for line

	for _, line := range strings.Split(a.CoverProfile, "\n") {
		profile := Profile{}
		tabs := strings.Split(line, "\t")
		if len(tabs) < 2 {
			continue
		}

		profile.Folder = strings.Trim(tabs[1], " ")
		//when the first one is ? => no coverage, no tests

		if strings.Trim(tabs[0], " ") == "?" {
			profile.Coverage = 0
			profile.Duration = 0

			a.profiles = append(a.profiles, profile)

			continue
		}

		fl := strings.TrimSuffix(tabs[2], "s")
		convertedValue, err := strconv.ParseFloat(fl, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		profile.Duration = convertedValue

		//get the coverage at a tring like coverage: 1.4% of statements in ./sdk/...
		reg := regexp.MustCompile(`coverage:\s+(\d+\.\d+)%`)
		matches := reg.FindStringSubmatch(tabs[3])
		if len(matches) != 2 {
			fmt.Println("could not find coverage")
			continue
		}

		//convert to a float
		convertedValue, err = strconv.ParseFloat(matches[1], 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		profile.Coverage = convertedValue

		a.profiles = append(a.profiles, profile)
	}
}
