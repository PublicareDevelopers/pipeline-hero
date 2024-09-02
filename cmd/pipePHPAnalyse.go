/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

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
	},
}

func init() {
	pipeCmd.AddCommand(pipePHPAnalyseCmd)
	pipePHPAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipePHPAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}

func analysePHPVulnCheck(analyser *code.PHPAnalyser, wg *sync.WaitGroup) {
	defer wg.Done()

	audit, err := cmds.GetComposerAudit()
	if err != nil {
		analyser.SetVulnCheck(audit)
		analyser.SetVulnCheckFail()

		//analyser.PushError("audit failed")
		//resp, err := sendVulnToPlatform(audit)
		//if err != nil {
		//	analyser.PushError("sending vuln to platform failed")
		//}
		//
		//color.White(fmt.Sprintf("Vuln Check: %s\n", resp))
	}
}
