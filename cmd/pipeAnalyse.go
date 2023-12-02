package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
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

		slackNotifySuccess(analyser, "go")
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
		wg.Done()
		return
	}

	analyser.SetGoVersion(version)

	color.Green("%s\n", version)
}

func analyseSDependencyGraph(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	dependencyGraph, err := cmds.GetDependencyGraph()
	if err != nil {
		analyser.PushError(fmt.Sprintf("cannot find the dependency graph: %s\n", err))
		color.Red("Error: %s\n", err)
		wg.Done()
		return
	}

	color.Green("setting dependency graph\n")
	analyser.SetDependencyGraph(dependencyGraph)

	toolchain, err := analyser.GetToolChainByDependencyGraph(dependencyGraph)
	if err != nil {
		color.Red("Error: %s\n", err)
		analyser.PushError(fmt.Sprintf("cannot find the toolchain: %s\n", err))
	}

	color.Green("%s\n", toolchain)

	//TODO => find only the ones by the own module
	dependencyUpdates := analyser.GetUpdatableDependencies()

	for _, depUpdate := range dependencyUpdates {
		color.Yellow("(used by %s) dependency update %s to %s\n", depUpdate.From, depUpdate.To, depUpdate.UpdateTo)
	}
}

func analyseTestCoverage(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	color.Green("\nrunning tests\n")
	out, err := exec.Command("go", "test", testSetup, fmt.Sprintf("-coverpkg=%s", testSetup), "-coverprofile=cover.cov").Output()
	if err != nil {
		fmt.Printf("%s\n", string(out))
		color.Red("Error: %s\n", err)
		analyser.PushError(fmt.Sprintf("Tests failed:\n%s\n", string(out)))
		wg.Done()
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
		wg.Done()
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
		wg.Done()
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
		wg.Done()
		return
	}

	fmt.Println(vulCheck)
	analyser.SetVulnCheck(vulCheck)
}
