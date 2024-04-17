package cmd

import (
	"fmt"
	"os"
	"pipeline-hero/sdk/platform"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pipeline-hero",
	Short: "some useful tools for your pipeline (alpha, not use in production expected you are a Publicare developer :) )",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sendVulnToPlatform(description string) (map[string]any, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic: sendVulnToPlatform")
		}
	}()

	repository := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	bitbucketProject := os.Getenv("BITBUCKET_PROJECT_KEY")

	platformClient := platform.New()

	platformClient.SetSecurityFixRequest(platform.SecurityFixRequest{
		Repository:       repository,
		BitbucketProject: bitbucketProject,
		Description:      description,
	})

	return platformClient.CreateSecurityTask()

}
