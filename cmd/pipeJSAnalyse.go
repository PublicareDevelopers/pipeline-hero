/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/fatih/color"
	"os"

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

		analyser := code.NewAnalyser()

		color.Green("analysing package.json\n")

		audit, err := cmds.GetNPMAudit()
		if err != nil {
			color.Red("%s\n\n", audit)
			color.Red("%s\n", err)
			analyser.PushError(audit)
			slackNotifyError(analyser, "npm audit failed")

			resp, err := sendVulnToPlatform("npm audit failed \n" + err.Error() + "\n" + audit)
			if err != nil {
				color.Red("error sending vuln to platform: %s\n", err)
				return
			}

			color.Green(fmt.Sprintf("vuln sent to platform: %+v\n", resp))

			os.Exit(255)
			return
		}

		fmt.Printf("%s\n", audit)

		slackNotifySuccess(analyser, "js")
	},
}

func init() {
	pipeCmd.AddCommand(pipeJSAnalyseCmd)
	pipeJSAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipeJSAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}
