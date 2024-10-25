package code

import (
	"encoding/json"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/npm"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/platform"
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

func (a *Analyser) GetDependenciesForPlatform(repository string) []platform.Dependency {
	out, err := cmds.AnalyseModule()
	if err != nil {
		a.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot find the module: %s\n", err))
		return nil
	}

	mod := GoMod{}
	err = json.Unmarshal([]byte(out), &mod)
	if err != nil {
		a.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot parse the module: %s\n", err))
		return nil
	}

	dependencies := make([]platform.Dependency, 0)
	for _, require := range mod.Require {
		dependencies = append(dependencies, platform.Dependency{
			Repository: repository,
			Name:       require.Path,
			Version:    require.Version,
			Language:   "go",
		})
	}

	return dependencies
}

func (a *JSAnalyser) GetDependenciesForPlatform(repository string) []platform.Dependency {
	out, err := cmds.AnalyseNPMModule()
	if err != nil {
		fmt.Println("having npm list error", err)
	}

	mod := npm.Mod{}
	err = json.Unmarshal([]byte(out), &mod)
	if err != nil {
		a.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot parse the module: %s\n", err))
		return nil
	}

	dependencies := make([]platform.Dependency, 0)
	for name, dep := range mod.Dependencies {
		dependencies = append(dependencies, platform.Dependency{
			Repository: repository,
			Name:       name,
			Version:    dep.Version,
			Language:   "js",
		})
	}

	return dependencies
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
// Deprecated: we now use the list of Requirements from the go.mod file: use here SetUpdatableDependencies
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

	//push them to the platform
	//platformClient := platform.New()
}
