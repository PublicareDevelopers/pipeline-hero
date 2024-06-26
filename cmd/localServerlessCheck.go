/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/qa"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var warnLevelMemory int

// localServerlessCheckCmd represents the localServerlessCheck command
var localServerlessCheckCmd = &cobra.Command{
	Use:   "check-serverless",
	Short: "checking the serverless configs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serverlessCheck, err := qa.ServerlessQA(rootDir, qa.WarnLevels{
			FunctionsMemorySize: warnLevelMemory,
		})

		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		for _, def := range serverlessCheck.Definitions {
			if len(def.Errors) > 0 {
				color.Red("%s", def.File)
				for _, e := range def.Errors {
					color.Red(e)
				}
			}

			if len(def.Warnings) > 0 {
				color.White("%s", def.File)
				for _, w := range def.Warnings {
					color.Yellow(w)
				}
			}

			if len(def.Errors) > 0 || len(def.Warnings) > 0 {
				color.White("-------------------")
			}
		}

		if len(serverlessCheck.MissingVars) > 0 {
			color.White("\nMissing environment variables:")
			for missingEnvVar, occurrences := range serverlessCheck.MissingVars {
				color.Yellow(fmt.Sprintf("%s | defined in: \n", missingEnvVar))
				for _, occurrence := range occurrences {
					color.Yellow(occurrence)
				}
				color.White("-------------------")
			}
		}
	},
}

func init() {
	localCmd.AddCommand(localServerlessCheckCmd)

	localServerlessCheckCmd.Flags().IntVarP(&warnLevelMemory, "warn-level-memory", "w", 4096, "Memory size warning level")
}
