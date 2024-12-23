/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

var exclude, sastPath string

// sastCmd represents the sast command
var sastCmd = &cobra.Command{
	Use:   "sast",
	Short: "make a SAST check of the current directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for key, value := range envVariables {
			err := os.Setenv(key, value)
			if err != nil {
				color.Red("Error: %s\n", err)
				os.Exit(255)
			}

			fmt.Println("set env variable", key, "to", value)
		}

		res, err := cmds.GoSASTProxy(exclude, sastPath)
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)

		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(sastCmd)

	sastCmd.Flags().StringVarP(&exclude, "exclude", "e", "G104", "exclude files from the SAST check")
	sastCmd.Flags().StringVarP(&sastPath, "path", "p", "./...", "path to the files to check")
}
