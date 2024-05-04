/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// localBuildCheckCmd represents the localBuildCheck command
var localBuildCheckCmd = &cobra.Command{
	Use:   "check-makefile",
	Short: "check the project config and local build status",
	Long:  `to can check, the project must follow some setup rules`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		res, err := cmds.CheckLocalBuild(rootDir)

		duration := time.Since(start)
		color.Green("local build duration: %v", duration)

		if err != nil {
			color.Red("Error: %v", err)
			color.Red(res)
			os.Exit(1)
		}

		color.Green("Makefiles are valid")
		color.Green(res)
	},
}

func init() {
	localCmd.AddCommand(localBuildCheckCmd)
}
