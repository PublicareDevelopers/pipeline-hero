package code

import "testing"

func TestAnalyser_GetToolChainByDependencyGraph(t *testing.T) {
	dependencyGraph := "github.com/PublicareDevelopers/pipeline-hero github.com/fatih/color@v1.15.0\n        github.com/PublicareDevelopers/pipeline-hero github.com/inconshreveable/mousetrap@v1.1.0\n        github.com/PublicareDevelopers/pipeline-hero github.com/mattn/go-colorable@v0.1.13\n        github.com/PublicareDevelopers/pipeline-hero github.com/mattn/go-isatty@v0.0.19\n        github.com/PublicareDevelopers/pipeline-hero github.com/spf13/cobra@v1.7.0\n        github.com/PublicareDevelopers/pipeline-hero github.com/spf13/pflag@v1.0.5\n        github.com/PublicareDevelopers/pipeline-hero golang.org/x/sys@v0.12.0\n        github.com/fatih/color@v1.15.0 github.com/mattn/go-colorable@v0.1.13\n        github.com/fatih/color@v1.15.0 github.com/mattn/go-isatty@v0.0.17\n        github.com/fatih/color@v1.15.0 golang.org/x/sys@v0.6.0\n        github.com/mattn/go-colorable@v0.1.13 github.com/mattn/go-isatty@v0.0.16\n        github.com/mattn/go-isatty@v0.0.19 golang.org/x/sys@v0.6.0\n        github.com/spf13/cobra@v1.7.0 github.com/cpuguy83/go-md2man/v2@v2.0.2\n        github.com/spf13/cobra@v1.7.0 github.com/inconshreveable/mousetrap@v1.1.0\n        github.com/spf13/cobra@v1.7.0 github.com/spf13/pflag@v1.0.5\n        github.com/spf13/cobra@v1.7.0 gopkg.in/yaml.v3@v3.0.1\n        github.com/mattn/go-isatty@v0.0.16 golang.org/x/sys@v0.0.0-20220811171246-fbc7d0a398ab\n        github.com/cpuguy83/go-md2man/v2@v2.0.2 github.com/russross/blackfriday/v2@v2.1.0\n        gopkg.in/yaml.v3@v3.0.1 gopkg.in/check.v1@v0.0.0-20161208181325-20d25e280405"

	a := &Analyser{}
	toolChain, err := a.GetToolChainByDependencyGraph(dependencyGraph)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	if toolChain == "" {
		t.Errorf("Error: toolChain is empty\n")
	}

	if toolChain != a.Toolchain {
		t.Errorf("Error: toolChain is not set\n")
	}
}
