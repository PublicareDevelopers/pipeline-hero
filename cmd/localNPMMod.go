/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"

	"github.com/spf13/cobra"
)

// localNPMModCmd represents the localNPMMod command
var localNPMModCmd = &cobra.Command{
	Use:   "npm-mod",
	Short: "output the npm module",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		a := code.JSAnalyser{}
		deps := a.GetDependenciesForPlatform("repository")
		for _, dep := range deps {
			cmd.Println(dep)
		}
	},
}

func init() {
	localCmd.AddCommand(localNPMModCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localNPMModCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localNPMModCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
