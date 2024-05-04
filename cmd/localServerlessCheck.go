/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/qa"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

// localServerlessCheckCmd represents the localServerlessCheck command
var localServerlessCheckCmd = &cobra.Command{
	Use:   "check-serverless",
	Short: "checking the serverless configs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		definitions, err := qa.ServerlessQA(rootDir)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		for _, def := range definitions {
			if len(def.Errors) > 0 {
				color.Red("%s", def.File)
				for _, e := range def.Errors {
					color.Red(e)
				}
			}

			if len(def.Warnings) > 0 {
				color.Yellow("%s", def.File)
				for _, w := range def.Warnings {
					color.Yellow(w)
				}
			}

			if len(def.Errors) > 0 || len(def.Warnings) > 0 {
				color.White("-------------------")
			}
		}
	},
}

func init() {
	localCmd.AddCommand(localServerlessCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localServerlessCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localServerlessCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
