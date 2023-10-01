package code

import (
	"regexp"
)

func (a *Analyser) GetToolChainByDependencyGraph(dependencyGraph string) (string, error) {
	reg := regexp.MustCompile(`(.*)toolchain(.*)`)
	matches := reg.FindStringSubmatch(dependencyGraph)

	if len(matches) > 0 {
		a.Toolchain = matches[1]
		return a.Toolchain, nil
	}

	a.Toolchain = "no toolchain found"

	return a.Toolchain, nil
}
