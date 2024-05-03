/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var rootDir string

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "", "root directory of the project")
	_ = localCmd.MarkPersistentFlagRequired("root")

}
