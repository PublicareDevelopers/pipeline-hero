/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// localScanCmd represents the localScan command
var localScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("scanning dependencies\n")
		dep, err := code.CheckDependencies()
		if err != nil {
			color.Red("Error: %s\n", err)
		}

		for _, depUpdate := range dep {
			color.Yellow("dependency update: %s\n", depUpdate)
		}

		color.Green("checking for vulnerabilities\n")
		var vulCheck string
		vulCheck, err = code.VulCheck()
		if err != nil {
			//not a crtical error
			color.Red("Error: %s\n", err)
		}

		fmt.Println(vulCheck)
	},
}

func init() {
	localCmd.AddCommand(localScanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localScanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localScanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
