package cmd

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/notifier"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
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

func slackNotifyError(analyser *code.Analyser, message string) {
	if !useSlack {
		return
	}

	color.Green("enabling Slack communication\n")
	handler, err := notifier.New("slack")
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = handler.Client.Validate()
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = handler.Client.BuildErrorBlocks(analyser, message)
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = handler.Client.Notify()
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}
}

func slackNotifySuccess(analyser *code.Analyser, pipeType string) {
	if !useSlack {
		return
	}

	color.Green("enabling Slack communication\n")
	handler, err := notifier.New("slack")
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = handler.Client.Validate()
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	switch pipeType {
	case "js":
		err = handler.Client.BuildJSBlocks(analyser)
	case "go":
		err = handler.Client.BuildBlocks(analyser)
	default:
		err = fmt.Errorf("unknown pipe type %s", pipeType)
	}

	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = handler.Client.Notify()
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}
}
