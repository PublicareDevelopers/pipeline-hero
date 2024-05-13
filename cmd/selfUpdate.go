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

// selfUpdateCmd represents the selfUpdate command
var selfUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "running an update for pipeline-hero",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.SelfUpdate()
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		if res != "" {
			color.Green(res)
			color.Green("updated pipeline-hero")
			return
		}

		color.Red("no update available")

	},
}

func init() {
	selfCmd.AddCommand(selfUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selfUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selfUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
