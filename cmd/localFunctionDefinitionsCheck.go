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

// localFunctionDefinitionsCheckCmd represents the localFunctionDefinitionsCheck command
var localFunctionDefinitionsCheckCmd = &cobra.Command{
	Use:   "check-functions",
	Short: "compare teststage and productivestage functions in yml files",
	Long:  `make sure that all functions defined in teststage are also defined in productivestage`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CheckFunctionDefinitions(rootDir)
		if err != nil {
			color.Red("Error: %v", err)
			color.Red(res)
			os.Exit(1)
		}

		color.Green("Function definitions are valid")
		color.Green(res)
	},
}

func init() {
	localCmd.AddCommand(localFunctionDefinitionsCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localFunctionDefinitionsCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localFunctionDefinitionsCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
