package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// pipeAnalyseCmd represents the pipeAnalyse command
var pipeAnalyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("go", "version").Output()
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		color.Green("%s\n", out)
	},
}

func init() {
	pipeCmd.AddCommand(pipeAnalyseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipeAnalyseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipeAnalyseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
