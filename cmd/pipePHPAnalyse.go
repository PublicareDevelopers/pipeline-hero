/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"pipeline-hero/sdk/cmds"
	"pipeline-hero/sdk/code"
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

		analyser := code.NewAnalyser()

		color.Green("analysing composer.json\n")

		audit, err := cmds.GetComposerAudit()
		if err != nil {
			color.Red("%s\n\n", audit)
			color.Red("%s\n", err)
			analyser.PushError(audit)
			slackNotifyError(analyser, "composer audit failed")

			os.Exit(255)
			return
		}

		fmt.Printf("%s\n", audit)

		slackNotifySuccess(analyser, "php")
	},
}

func init() {
	pipeCmd.AddCommand(pipePHPAnalyseCmd)
	pipePHPAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipePHPAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}
