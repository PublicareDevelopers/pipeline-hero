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
}
