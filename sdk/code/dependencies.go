package code

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/fatih/color"
	"strings"
)

var maxDependencyChecks = 100

func (a *Analyser) SetUpdatableRequirements() {
	updatable := make([]RequireUpdate, 0)
	for _, require := range a.Module.Require {
		update, err := cmds.GetUpdateVersion(require.Path)
		if err != nil {
			fmt.Printf("error getting update version for %s: %s\n", require.Path, err)
			continue
		}

		if update != "" {
			updatable = append(updatable, RequireUpdate{
				Path:             require.Path,
				Version:          require.Version,
				AvailableVersion: update,
				Indirect:         require.Indirect,
			})
		}
	}

	for _, depUpdate := range updatable {
		color.Yellow("dependency update available for  %s: %s -> %s\n", depUpdate.Path, depUpdate.Version, depUpdate.AvailableVersion)
	}

	a.lock.Lock()
	a.Updates = updatable
	a.lock.Unlock()
}

// GetUpdatableDependencies
// Deprecated: we now use the list of Requirements from the go.mod file: use here SetUpdatableDependencies
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

// parseDependencyGraph
// Deprecated: we now use the list of Requirements from the go.mod file
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

		a.lock.Lock()
		a.dependencies = append(a.dependencies, Dependency{
			From:      original,
			To:        dependency,
			Updatable: updatable,
			UpdateTo:  updateTo,
		})
		a.lock.Unlock()
	}

	/*if len(a.dependencies) > maxDependencyChecks {
		a.PushWarning(
			fmt.Sprintf("Only the first %d dependencies are checked for updates. Have a total of %d", maxDependencyChecks, len(a.dependencies)))
	}*/
}
