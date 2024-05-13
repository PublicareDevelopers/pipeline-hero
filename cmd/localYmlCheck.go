/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

// localYmlCheckCmd represents the localYmlCheck command
var localYmlCheckCmd = &cobra.Command{
	Use:   "check-yml",
	Short: "check the serverless yml files searching for not available binaries in bin/ folder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CheckYml(rootDir)
		if err != nil {
			color.Red("Error: %v", err)
			color.Red(res)
			os.Exit(1)
		}

		color.Green("YML files are valid")
		color.Green(res)
	},
}

func init() {
	localCmd.AddCommand(localYmlCheckCmd)
}
