package cmds

import (
	"strings"
	"testing"
)

func TestGetDependencyGraph(t *testing.T) {
	graph, err := GetDependencyGraph()
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	if graph == "" {
		t.Errorf("Error: graph is empty\n")
	}

	if !strings.Contains(graph, "github.com") {
		t.Errorf("Error: graph is not correct\n")
	}
}
