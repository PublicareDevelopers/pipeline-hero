package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/slack"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"sync"
)

const BAD = 50.0
const MEDIUM = 75.0

var testSetup string
var envVariables map[string]string

// pipeAnalyseCmd represents the pipeAnalyse command
var pipeAnalyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for key, value := range envVariables {
			err := os.Setenv(key, value)
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			fmt.Println("set env variable", key, "to", value)
		}

		analyser := code.NewAnalyser().SetThreshold(coverageThreshold)

		var wg sync.WaitGroup

		wg.Add(1)
		go analyseVersion(analyser, &wg)

		wg.Add(1)
		go analyseSDependencyGraph(analyser, &wg)

		wg.Add(1)
		go analyseTestCoverage(analyser, &wg)

		wg.Add(1)
		go analyseVulnCheck(analyser, &wg)

		wg.Wait()

		if !useSlack {
			if len(analyser.GetErrors()) > 0 {
				color.Red("pipeline-hero failed\n")
				os.Exit(255)
			}
			return
		}

		client, err := slack.NewClient()
		if err != nil {
			color.Red("error at slack handling: %s\n", err)
			os.Exit(255)
		}

		err = client.StartConversation(analyser, "go")
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		if len(analyser.GetErrors()) > 0 {
			color.Red("pipeline-hero failed\n")
			os.Exit(255)
		}
	},
}

func init() {
	pipeCmd.AddCommand(pipeAnalyseCmd)
	pipeAnalyseCmd.Flags().StringVarP(&testSetup, "test-setup", "t", "./...", "Test setup to use")
	pipeAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipeAnalyseCmd.Flags().Float64VarP(&coverageThreshold, "coverage-threshold", "c", 75.0, "Coverage threshold to use")
	pipeAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}

func analyseVersion(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	version, err := cmds.GetVersion()
	if err != nil {
		color.Red("Error: %s\n", err)
		analyser.PushError(fmt.Sprintf("cannot find the go version: %s\n", err))
		return
	}

	analyser.SetGoVersion(version)

	color.Green("%s\n", version)
}

func analyseSDependencyGraph(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	analyser.SetModule()
	analyser.SetUpdatableRequirements()

	dependencyGraph, err := cmds.GetDependencyGraph()
	if err != nil {
		analyser.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot find the dependency graph: %s\n", err))
		color.Red("Error: %s\n", err)
		return
	}

	color.Green("setting dependency graph\n")

	analyser.SetDependencyGraph(dependencyGraph)
}

func analyseTestCoverage(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	color.Green("\nrunning tests\n")
	out, err := exec.Command("go", "test", testSetup, fmt.Sprintf("-coverpkg=%s", testSetup), "-coverprofile=cover.cov").Output()
	if err != nil {
		color.Red("Error: %s\n", err)
		analyser.TestResult = string(out)
		color.Red("Tests failed:\n%s\n", string(out))
		analyser.PushError(fmt.Sprintf("Tests failed:\n%sThe results can be found in the thread or have a look at the pipeline\n", err))
		return
	}

	coverProfile := string(out)
	analyser.SetCoverProfile(coverProfile)

	fmt.Println(coverProfile)
	color.Green("running coverage\n")

	//go tool cover -func=cover.cov | tee coverage.txt
	out, err = exec.Command("go", "tool", "cover", "-func=cover.cov").Output()
	if err != nil {
		fmt.Printf("%s\n", string(out))
		color.Red("Error: %s\n", err)
		analyser.PushError(fmt.Sprintf("Coverage failed:\n%s\n", string(out)))
		return
	}

	//we have something like total:  (statements)    0.0%
	//grep the total amount with a regex
	totalText := string(out)

	analyser.SetCoverageByTotal(totalText)

	fmt.Println(analyser.GetCoverageInterpretation())

	err = analyser.CheckThreshold()
	if err != nil {
		color.Red("%s\n", err)
		analyser.PushError(fmt.Sprintf("coverage threshold not met: have  %.2f  percent", analyser.Coverage))
		return
	}
}

func analyseVulnCheck(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	color.Green("starting vuln check\n")

	var vulCheck string
	vulCheck, err := cmds.VulnCheck(testSetup)
	if err != nil {
		color.Red("%s\n", err)
		analyser.PushError(err.Error())
		return
	}

	fmt.Println(vulCheck)
	analyser.SetVulnCheck(vulCheck)
}
