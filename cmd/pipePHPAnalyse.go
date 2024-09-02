/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/npm"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/php"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/slack"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

var testFolders []string
var phpUnitCmd string

// pipePHPAnalyseCmd represents the pipePHPAnalyse command
var pipePHPAnalyseCmd = &cobra.Command{
	Use:   "php-analyse",
	Short: "check some php components; make sure composer is installed at the docker image",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("setting environment if some arguments are given\n")
		for key, value := range envVariables {
			err := os.Setenv(key, value)
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			fmt.Println("set env variable", key, "to", value)
		}

		analyser := code.NewPHPAnalyser().SetThreshold(coverageThreshold)

		var wg sync.WaitGroup

		wg.Add(1)
		go analysePHPVulnCheck(analyser, &wg)

		wg.Add(1)
		go runPHPUnitTests(analyser, &wg)

		wg.Add(1)
		go analysePHPOutDates(analyser, &wg)

		wg.Wait()

		if !useSlack {
			if len(analyser.GetErrors()) > 0 {
				color.Red("pipeline-hero failed, have %d errors", len(analyser.GetErrors()))
			}

			if analyser.HasVulnCheckFail {
				color.Red("Vuln Check failed")
				color.Red(analyser.VulnCheck)
			}

			if analyser.HasWarnings {
				color.Yellow("Warnings found\n")
				for _, warning := range analyser.Warnings {
					color.Yellow(warning)
				}
			}

			os.Exit(255)
			return
		}

		client, err := slack.NewClient()
		if err != nil {
			color.Red("error at slack handling: %s\n", err)
			os.Exit(255)
		}

		err = client.StartPHPConversation(analyser)
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
	pipeCmd.AddCommand(pipePHPAnalyseCmd)
	pipePHPAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipePHPAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
	pipePHPAnalyseCmd.Flags().Float64VarP(&coverageThreshold, "coverage-threshold", "c", 75.0, "Coverage threshold to use")
	pipePHPAnalyseCmd.Flags().StringSliceVarP(&testFolders, "test-folders", "t", []string{}, "Test folders to use")
	pipePHPAnalyseCmd.Flags().StringVarP(&phpUnitCmd, "phpunit-cmd", "p", "./vendor/bin/phpunit", "PHPUnit command to use")
}

func analysePHPVulnCheck(analyser *code.PHPAnalyser, wg *sync.WaitGroup) {
	defer wg.Done()

	audit, err := cmds.GetComposerAudit()
	if err != nil {
		analyser.SetVulnCheck(audit)
		analyser.SetVulnCheckFail()

		analyser.PushError("audit failed")
		resp, err := sendVulnToPlatform(audit)
		if err != nil {
			analyser.PushError("sending vuln to platform failed")
		}

		color.White(fmt.Sprintf("Vuln Check: %s\n", resp))
	}
}

func runPHPUnitTests(analyser *code.PHPAnalyser, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, folder := range testFolders {
		out, err := cmds.RunPHPUnitTest(phpUnitCmd, folder)
		if err != nil {
			analyser.PushError(fmt.Sprintf("test failed at %s: %s\n", folder, out))
			continue
		}

		color.Green("successfully tested: %s\n", folder)
	}
}

func analysePHPOutDates(analyser *code.PHPAnalyser, wg *sync.WaitGroup) {
	defer wg.Done()

	outdatedJson, err := cmds.GetComposerOutDates()
	if err != nil {
		color.Red("Error getting OutDates: %s\n", err)
		return
	}

	type OutDatesJson struct {
		Installed []php.OutDate `json:"installed"`
	}

	var outDatesRaw OutDatesJson
	err = json.Unmarshal([]byte(outdatedJson), &outDatesRaw)

	var outDates []php.OutDate

	for _, outDateEl := range outDatesRaw.Installed {
		outDate, err := outDateEl.ParseOutDate()
		if err != nil {
			color.Red("Error: %s\n", err)
			continue
		}

		outDate.Rate()

		outDates = append(outDates, *outDate)

		if outDate.Rating.StatusCode == npm.OutDateRatingNewMajorVersionAvailable {
			analyser.PushWarning(fmt.Sprintf("%s: %s", outDate.Name, outDate.Rating.Message))
		}
	}

	analyser.SetOutDates(outDates)
}
