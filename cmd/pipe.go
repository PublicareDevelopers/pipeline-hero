package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"pipeline-hero/sdk/code"
	"pipeline-hero/sdk/slack"
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

	client, err := slack.NewClient()
	if err != nil {
		color.Red("error at slack handling: %s\n", err)
		os.Exit(255)
	}

	err = client.BuildErrorBlocks(analyser, message)
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = client.SendProgressSlackBlocks(client.Blocks)
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}
}

func slackNotifySuccess(analyser *code.Analyser, pipeType string) {
	if !useSlack {
		return
	}

	client, err := slack.NewClient()
	if err != nil {
		color.Red("error at slack handling: %s\n", err)
		os.Exit(255)
	}

	switch pipeType {
	case "js":
		err = client.BuildJSBlocks(analyser)
	case "go":
		err = client.BuildBlocks(analyser)
	case "php":
		err = client.BuildPHPBlocks(analyser)
	default:
		err = fmt.Errorf("unknown pipe type %s", pipeType)
	}

	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}

	err = client.SendProgressSlackBlocks(client.Blocks)
	if err != nil {
		color.Red("Error: %s\n", err)
		os.Exit(255)
	}
}
