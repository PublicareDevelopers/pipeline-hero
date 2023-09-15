package code

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

func Analyse() ([]string, error) {
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
