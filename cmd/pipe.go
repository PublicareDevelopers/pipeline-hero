package cmd

import (
	"github.com/spf13/cobra"
)

var useSlack bool
var coverageThreshold float64

// pipeCmd represents the pipe command
var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(pipeCmd)
	pipeCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipeCmd.Flags().Float64VarP(&coverageThreshold, "coverage-threshold", "c", 75.0, "Coverage threshold to use")
}
