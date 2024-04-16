package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"pipeline-hero/sdk/cmds"
	"pipeline-hero/sdk/code"

	"github.com/spf13/cobra"
)

// localDependencyScanCmd represents the localDependencyScan command
var localDependencyScanCmd = &cobra.Command{
	Use:   "scan-dependencies",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		analyser := code.NewAnalyser()

		dependencyGraph, err := cmds.GetDependencyGraph()
		if err != nil {
			color.Red("Error: %s\n", err)
			os.Exit(255)
		}

		color.Green("setting dependency graph\n")
		analyser.SetDependencyGraph(dependencyGraph)

		for _, dependency := range analyser.GetDependencyGraph() {
			fmt.Printf("%s -> %s\n", dependency.From, dependency.To)
		}

		color.Green("have %d dependencies\n", len(analyser.GetDependencyGraph()))
	},
}

func init() {
	localCmd.AddCommand(localDependencyScanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localDependencyScanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localDependencyScanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
