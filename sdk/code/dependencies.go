package code

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"strings"
)

var maxDependencyChecks = 100

func (a *Analyser) GetUpdatableDependencies() []Dependency {
	updatable := make([]Dependency, 0)
	for _, dependency := range a.dependencies {
		if dependency.Updatable {
			updatable = append(updatable, dependency)
		}
	}

	return updatable
}

func (a *Analyser) GetDependencyGraph() []Dependency {
	return a.dependencies
}

func (a *Analyser) parseDependencyGraph() {
	for count, line := range strings.Split(a.DependencyGraph, "\n") {
		if len(line) == 0 {
			continue
		}
		line = strings.Trim(line, " ")

		words := strings.Split(line, " ")
		original := words[0]
		dependency := words[1]

		updatable := false
		updateTo := ""

		if count < maxDependencyChecks {
			update, err := cmds.GetUpdateVersion(dependency)
			if err == nil {
				updatable = update != ""
				updateTo = update
			}
		}

		a.dependencies = append(a.dependencies, Dependency{
			From:      original,
			To:        dependency,
			Updatable: updatable,
			UpdateTo:  updateTo,
		})
	}

	if len(a.dependencies) > maxDependencyChecks {
		a.PushWarning(
			fmt.Sprintf("Only the first %d dependencies are checked for updates. Have a total of %d", maxDependencyChecks, len(a.dependencies)))
	}
}
