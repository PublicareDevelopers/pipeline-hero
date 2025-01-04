/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/spf13/cobra"
	"os"
)

var codereviewVueLevel, codereviewVueOutput string

// codereviewVueCmd represents the codereviewVue command
var codereviewVueCmd = &cobra.Command{
	Use:   "vue",
	Short: "making a code review of your vue code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		codereview, err := cmds.CodeReviewVueMess(codereviewVueLevel, codereviewVueOutput)
		if err != nil {
			cmd.Println("pipeline-hero failed", err, codereview)
			os.Exit(1)
		}
		cmd.Println(codereview)
	},
}

func init() {
	codereviewCmd.AddCommand(codereviewVueCmd)

	codereviewVueCmd.Flags().StringVarP(&codereviewVueLevel, "level", "l", "all", "level of the code review; default is all; can be set to error only")
	codereviewVueCmd.Flags().StringVarP(&codereviewVueOutput, "output", "o", "table", "output format; default is table; can be set to table, json or text")
}
