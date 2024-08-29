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
	"github.com/PublicareDevelopers/pipeline-hero/sdk/slack"
	"github.com/fatih/color"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// pipeJSAnalyseCmd represents the pipeJSAnalyse command
var pipeJSAnalyseCmd = &cobra.Command{
	Use:   "js-analyse",
	Short: "check some frontend components; make sure you are in the same directory where package.json is located",
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

		analyser := code.NewJSAnalyser().SetThreshold(coverageThreshold)

		var wg sync.WaitGroup

		wg.Add(1)
		go analyseJSVulnCheck(analyser, &wg)

		wg.Add(1)
		go analyseJSOutDates(analyser, &wg)

		wg.Wait()

		if !useSlack {
			if len(analyser.GetErrors()) > 0 {
				color.Red("pipeline-hero failed")
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

		err = client.StartJSConversation(analyser)
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
	pipeCmd.AddCommand(pipeJSAnalyseCmd)
	pipeJSAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipeJSAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}

func analyseJSVulnCheck(analyser *code.JSAnalyser, wg *sync.WaitGroup) {
	defer wg.Done()

	audit, err := cmds.GetNPMAudit()
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

func analyseJSOutDates(analyser *code.JSAnalyser, wg *sync.WaitGroup) {
	defer wg.Done()

	outdatedJson, err := cmds.GetNPMOutdated()
	if err != nil {
		color.Red("Error getting OutDates: %s\n", err)
		return
	}

	var outDatesRaw map[string]npm.OutDate
	err = json.Unmarshal([]byte(outdatedJson), &outDatesRaw)

	var outDates []npm.OutDate

	for name, outDateEl := range outDatesRaw {
		outDate, err := npm.NewOutDate(name, outDateEl.Current, outDateEl.Wanted, outDateEl.Latest, outDateEl.Dependent)
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
