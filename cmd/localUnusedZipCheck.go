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

// localUnusedZipCheckCmd represents the localUnusedZipCheck command
var localUnusedZipCheckCmd = &cobra.Command{
	Use:   "check-zip",
	Short: "searching for unused zips",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CheckUnusedZips(rootDir)
		if err != nil {
			color.Red("Error: %v", err)
			color.Red(res)
			os.Exit(1)
		}

		color.Green("no unused zips found")
		color.Green(res)
	},
}

func init() {
	localCmd.AddCommand(localUnusedZipCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localUnusedZipCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localUnusedZipCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
