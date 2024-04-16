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
)

// localScanCmd represents the localScan command
var localScanCmd = &cobra.Command{
	Use:   "scan-vuln",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		vulCheck, err := cmds.VulnCheck(testSetup)
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
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
