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

var codePath string

// codereviewGoCmd represents the codereviewGo command
var codereviewGoCmd = &cobra.Command{
	Use:   "go",
	Short: "make a codereview of your go code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := cmds.CodeReview(codePath)
		if err != nil {
			color.Red("pipeline-hero failed\n", err)
		}

		color.Green("code review result: ")
		fmt.Println(res)
	},
}

func init() {
	codereviewCmd.AddCommand(codereviewGoCmd)

	codereviewGoCmd.Flags().StringVarP(&codePath, "path", "p", "", "path to the go code")
	_ = codereviewGoCmd.MarkFlagRequired("path")
}
