package code

import (
	"regexp"
)

func (a *Analyser) GetToolChainByDependencyGraph(dependencyGraph string) (string, error) {
	reg := regexp.MustCompile(`(.*)toolchain(.*)`)
	matches := reg.FindStringSubmatch(dependencyGraph)

	if len(matches) > 0 {
		return matches[1], nil
	}

	return "no toolchain found", nil
}
