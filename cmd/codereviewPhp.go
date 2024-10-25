/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var phpStanLevel, phpPath string

// codereviewPhpCmd represents the codereviewPhp command
var codereviewPhpCmd = &cobra.Command{
	Use:   "php",
	Short: "make a codereview of your php code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CodeReviewPHPStan(phpPath, phpStanLevel)
		if err != nil {
			color.Red("pipeline-hero failed\n", err)
		}

		color.Green("code review result: ")
		fmt.Println(res)
	},
}

func init() {
	codereviewCmd.AddCommand(codereviewPhpCmd)

	codereviewPhpCmd.Flags().StringVarP(&phpStanLevel, "level", "l", "5", "phpstan level")
	codereviewPhpCmd.Flags().StringVarP(&phpPath, "path", "p", "./...", "path to the php code")
}
