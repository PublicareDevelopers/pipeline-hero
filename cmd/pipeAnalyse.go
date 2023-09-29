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
)

const BAD = 50.0
const MEDIUM = 75.0

var useSlack bool
var testSetup string
var coverageThreshold float64
var envVariables map[string]string

// pipeAnalyseCmd represents the pipeAnalyse command
var pipeAnalyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := slack.New()
		analyser := code.NewAnalyser().SetThreshold(coverageThreshold)

		version, err := cmds.GetVersion()
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		color.Green("%s\n", version)
		client.GoVersion = version

		dependencyGraph, err := cmds.GetDependencyGraph()
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		toolchain, err := analyser.GetToolChainByDependencyGraph(dependencyGraph)
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		client.GoToolchainVersion = toolchain
		color.Green("%s\n", toolchain)

		color.Green("setting environment if some arguments are given\n")
		for key, value := range envVariables {
			err := os.Setenv(key, value)
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			fmt.Println("set env variable", key, "to", value)
		}

		color.Green("\nrunning tests\n")
		out, err := exec.Command("go", "test", testSetup, fmt.Sprintf("-coverpkg=%s", testSetup), "-coverprofile=cover.cov").Output()
		if err != nil {
			fmt.Printf("%s\n", string(out))
			color.Red("Error: %s\n", err)
			os.Exit(255)
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
			os.Exit(255)
		}

		//we have something like total:  (statements)    0.0%
		//grep the total amount with a regex
		totalText := string(out)

		analyser.SetCoverageByTotal(totalText)

		fmt.Println(analyser.GetCoverageInterpretation())

		dep, err := code.CheckDependencies()
		if err != nil {
			//not a crtical error
			client.Errors = append(client.Errors, err.Error())
		}

		client.DependencyUpdates = dep

		for _, depUpdate := range dep {
			color.Yellow("dependency update: %s\n", depUpdate)
		}

		var vulCheck string
		vulCheck, err = code.VulCheck()
		if err != nil {
			//not a crtical error
			client.Errors = append(client.Errors, err.Error())
		}

		client.VulCheck = vulCheck

		if useSlack {
			color.Green("enabling Slack communication\n")
			client, err = client.InitConfiguration()
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			err = client.BuildBlocksByBitbucket(fmt.Sprintf("have %.2f  percent coverage\n", analyser.Coverage)).
				NotifyWithBlocks()

			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}
		}
	},
}

func init() {
	pipeCmd.AddCommand(pipeAnalyseCmd)

	pipeAnalyseCmd.Flags().Float64VarP(&coverageThreshold, "coverage-threshold", "c", 75.0, "Coverage threshold to use")
	pipeAnalyseCmd.Flags().StringVarP(&testSetup, "test-setup", "t", "./...", "Test setup to use")
	pipeAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipeAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}
