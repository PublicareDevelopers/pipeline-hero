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

// localAllCheckCmd represents the localAllCheck command
var localAllCheckCmd = &cobra.Command{
	Use:   "check-complete",
	Short: "check the makefiles and the pipeline yml files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CheckBuild(rootDir)
		if err != nil {
			color.Red("Error: %v", err)
			color.Red(res)
			os.Exit(1)
		}

		color.Green("Makefiles and YML files are valid")
		color.Green(res)
	},
}

func init() {
	localCmd.AddCommand(localAllCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localAllCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localAllCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
