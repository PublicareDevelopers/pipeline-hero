package code

import (
	"strings"
	"testing"
)

func TestAnalyser_SetDependencyGraph(t *testing.T) {
	dp := "/pipeline-hero github.com/fatih/color@v1.15.0\n        /pipeline-hero github.com/inconshreveable/mousetrap@v1.1.0\n        /pipeline-hero github.com/mattn/go-colorable@v0.1.13\n        /pipeline-hero github.com/mattn/go-isatty@v0.0.19\n        /pipeline-hero github.com/spf13/cobra@v1.7.0\n        /pipeline-hero github.com/spf13/pflag@v1.0.5\n        /pipeline-hero golang.org/x/sys@v0.12.0\n        github.com/fatih/color@v1.15.0 github.com/mattn/go-colorable@v0.1.13\n        github.com/fatih/color@v1.15.0 github.com/mattn/go-isatty@v0.0.17\n        github.com/fatih/color@v1.15.0 golang.org/x/sys@v0.6.0\n        github.com/mattn/go-colorable@v0.1.13 github.com/mattn/go-isatty@v0.0.16\n        github.com/mattn/go-isatty@v0.0.19 golang.org/x/sys@v0.6.0\n        github.com/spf13/cobra@v1.7.0 github.com/cpuguy83/go-md2man/v2@v2.0.2\n        github.com/spf13/cobra@v1.7.0 github.com/inconshreveable/mousetrap@v1.1.0\n        github.com/spf13/cobra@v1.7.0 github.com/spf13/pflag@v1.0.5\n        github.com/spf13/cobra@v1.7.0 gopkg.in/yaml.v3@v3.0.1\n        github.com/mattn/go-isatty@v0.0.16 golang.org/x/sys@v0.0.0-20220811171246-fbc7d0a398ab\n        github.com/cpuguy83/go-md2man/v2@v2.0.2 github.com/russross/blackfriday/v2@v2.1.0\n        gopkg.in/yaml.v3@v3.0.1 gopkg.in/check.v1@v0.0.0-20161208181325-20d25e280405"
	lines := strings.Split(dp, "\n")

	a := NewAnalyser().SetDependencyGraph(dp)
	if a.DependencyGraph != dp {
		t.Errorf("SetDependencyGraph() = %v, want %v", a.DependencyGraph, dp)
	}

	if len(a.dependencies) != len(lines) {
		t.Errorf("SetDependencyGraph() = %v, want %v", len(a.dependencies), len(lines))
	}

	dependency := a.dependencies[0]
	if dependency.From != "/pipeline-hero" {
		t.Errorf("SetDependencyGraph() = %v, want %v", dependency.From, "/pipeline-hero")
	}

	if dependency.To != "github.com/fatih/color@v1.15.0" {
		t.Errorf("SetDependencyGraph() = %v, want %v", dependency.To, "github.com/fatih/color@v1.15.0")
	}

	dependency = a.dependencies[11]
	if dependency.From != "github.com/mattn/go-isatty@v0.0.19" {
		t.Errorf("SetDependencyGraph() = %v, want %v", dependency.From, "github.com/mattn/go-isatty@v0.0.19")
	}

	if dependency.To != "golang.org/x/sys@v0.6.0" {
		t.Errorf("SetDependencyGraph() = %v, want %v", dependency.To, "golang.org/x/sys@v0.6.0")
	}

	for _, dep := range a.dependencies {
		if dep.Updatable {
			if dep.UpdateTo == "" {
				t.Errorf("SetDependencyGraph() = %v, want a value", dep.UpdateTo)
			}

			if strings.Contains(dep.UpdateTo, "[") {
				t.Errorf("SetDependencyGraph() = %v, want a value without [", dep.UpdateTo)
			}

			if strings.Contains(dep.UpdateTo, "]") {
				t.Errorf("SetDependencyGraph() = %v, want a value without ]", dep.UpdateTo)
			}
		}
	}

}

func TestAnalyser_GetUpdatableDependencies(t *testing.T) {
	a := NewAnalyser()

	a.dependencies = append(a.dependencies, Dependency{
		From:      "/pipeline-hero",
		To:        "github.com/fatih/color@v1.15.0",
		Updatable: true,
		UpdateTo:  "v1.15.1",
	})

	a.dependencies = append(a.dependencies, Dependency{
		From:      "/pipeline-hero",
		To:        "github.com/inconshreveable/mousetrap@v1.1.0",
		Updatable: false,
		UpdateTo:  "",
	})

	updateable := a.GetUpdatableDependencies()
	if len(updateable) != 1 {
		t.Errorf("GetUpdatableDependencies() = %v, want %v", len(updateable), 1)
	}
}

func TestAnalyser_GetDependencyGraph(t *testing.T) {
	a := NewAnalyser()

	a.dependencies = append(a.dependencies, Dependency{
		From:      "/pipeline-hero",
		To:        "github.com/fatih/color@v1.15.0",
		Updatable: true,
		UpdateTo:  "v1.15.1",
	})

	a.dependencies = append(a.dependencies, Dependency{
		From:      "/pipeline-hero",
		To:        "github.com/inconshreveable/mousetrap@v1.1.0",
		Updatable: false,
		UpdateTo:  "",
	})

	deps := a.GetDependencyGraph()
	if len(deps) != 2 {
		t.Errorf("GetUpdatableDependencies() = %v, want %v", len(deps), 1)
	}
}

func TestAnalyser_GetDependenciesForPlatform(t *testing.T) {
	a := NewAnalyser()
	deps := a.GetDependenciesForPlatform("test")
	if deps == nil {
		t.Fatalf("GetDependenciesForPlatform() = %v, want > %v", deps, nil)
	}

	if len(deps) == 0 {
		t.Fatalf("GetDependenciesForPlatform() = %v, want > %v", len(deps), 0)
	}
}
