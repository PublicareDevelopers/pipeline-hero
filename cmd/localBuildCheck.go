/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

// localBuildCheckCmd represents the localBuildCheck command
var localBuildCheckCmd = &cobra.Command{
	Use:   "check-makefile",
	Short: "check the project config and local build status",
	Long:  `to can check, the project must follow some setup rules`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CheckLocalBuild(rootDir)
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
