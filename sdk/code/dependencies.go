package code

import (
	"bytes"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"os/exec"
	"regexp"
	"strings"
)

func (a *Analyser) GetUpdatableDependencies() []Dependency {
	updatable := make([]Dependency, 0)
	for _, dependency := range a.dependencies {
		if dependency.Updatable {
			updatable = append(updatable, dependency)
		}
	}

	return updatable
}

func (a *Analyser) parseDependencyGraph() {
	for _, line := range strings.Split(a.DependencyGraph, "\n") {
		if len(line) == 0 {
			continue
		}
		line = strings.Trim(line, " ")

		words := strings.Split(line, " ")
		original := words[0]
		dependency := words[1]

		updatable := false
		updateTo := ""

		update, err := cmds.GetUpdateVersion(dependency)
		if err == nil {
			updatable = update != ""
			updateTo = update
		}

		a.dependencies = append(a.dependencies, Dependency{
			From:      original,
			To:        dependency,
			Updatable: updatable,
			UpdateTo:  updateTo,
		})
	}
}

func CheckDependencies() ([]string, error) {
	scanned := make(map[string]bool)

	graph, err := exec.Command("go", "mod", "graph").Output()
	if err != nil {
		return nil, err
	}

	var dependencyUpdates []string
	//reading graph line for line
	lines := bytes.Split(graph, []byte("\n"))
	for _, line := range lines {
		//if it is an empty line continue
		if len(line) == 0 {
			continue
		}
		//split the line by empty space
		words := bytes.Split(line, []byte(" "))
		original := string(words[0])

		if _, ok := scanned[original]; ok {
			continue
		}

		scanned[original] = true

		out, err := exec.Command("go", "list", "-m", "-u", original).Output()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		//check if out have something like []
		reg := regexp.MustCompile(`\[(.*)\]`)
		matches := reg.FindStringSubmatch(string(out))
		if len(matches) > 0 {
			dependencyUpdates = append(dependencyUpdates, fmt.Sprintf("%s: %s", original, matches[1]))
		}

	}

	return dependencyUpdates, nil
}
