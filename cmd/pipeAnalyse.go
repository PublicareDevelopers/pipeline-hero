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
	"regexp"
	"strconv"
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

		//find complete line where toolchain is mentioned
		reg := regexp.MustCompile(`(.*)toolchain(.*)`)
		matches := reg.FindStringSubmatch(dependencyGraph)
		if len(matches) > 0 {
			client.GoToolchainVersion = matches[1]
			color.Green("%s\n", matches[1])

		}

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

		fmt.Println(string(out))
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

		reg = regexp.MustCompile(`total:\s+\((\w+)\)\s+(\d+\.\d+)%`)
		matches = reg.FindStringSubmatch(totalText)
		if len(matches) != 3 {
			color.Red("Error: could not find total coverage\n have %s\n", totalText)
			os.Exit(255)
		}

		//convert to a float
		total, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		if total < BAD {
			color.Red("coverage is BAD, have %.2f  percent\n", total)
		}

		if total < MEDIUM && total > BAD {
			color.Yellow("coverage is ok, have %.2f  percent\n", total)
		}

		if total >= MEDIUM {
			color.Green("coverage is good, have %.2f  percent\n", total)
		}

		if total < coverageThreshold {
			color.Red("coverage threshold %.2f  not reached, have %.2f \n", coverageThreshold, total)
			os.Exit(255)
		}

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

			err = client.BuildBlocksByBitbucket(fmt.Sprintf("have %.2f  percent coverage\n", total)).
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
