/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/notifier"
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
		analyser := code.NewAnalyser().SetThreshold(coverageThreshold)

		color.Green("analysing package.json\n")

		audit, err := cmds.GetNPMAudit()
		if err != nil {
			color.Red("%s\n\n", audit)
			color.Red("%s\n", err)
			os.Exit(255)
			return
		}

		fmt.Printf("%s\n", audit)

		if useSlack {
			color.Green("enabling Slack communication\n")
			handler, err := notifier.New("slack")
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			err = handler.Client.Validate()
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			err = handler.Client.BuildBlocks(analyser)
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			err = handler.Client.Notify()
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}
		}
	},
}

func init() {
	pipeCmd.AddCommand(pipeJSAnalyseCmd)
}
