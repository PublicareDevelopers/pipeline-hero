package cmd

import (
	"github.com/spf13/cobra"
)

// codereviewCmd represents the codereview command
var codereviewCmd = &cobra.Command{
	Use:   "codereview",
	Short: "start a codereview",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(codereviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// codereviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// codereviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
