package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (slack *Slack) Notify(message string) error {
	reader := bytes.NewReader([]byte(fmt.Sprintf(
		`{"text":"%s"}`,
		fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC822), message),
	)))

	req, err := http.NewRequest(http.MethodPost, slack.WebhookURL, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing slack response body: %s\n", err.Error())
		}
	}()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Slack response status code: %d", resp.StatusCode))
	}

	return nil
}

func (slack *Slack) NotifyWithBlocks() error {
	blockJson, err := json.Marshal(map[string]any{
		"blocks": slack.Blocks,
	})
	if err != nil {
		return err
	}

	reader := bytes.NewReader(blockJson)

	req, err := http.NewRequest(http.MethodPost, slack.WebhookURL, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing slack response body: %s\n", err.Error())
		}
	}()

	if resp.StatusCode != 200 {
		//return errors.New(fmt.Sprintf("Slack response status code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(blockJson))
	fmt.Printf("Slack response: %+v\n", string(body))

	return nil
}

func BuildBitBucketMessage(message string) string {
	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	bitbucketUrl := fmt.Sprintf("%s/addon/pipelines/home#!/results/%s", origin, commit)

	return fmt.Sprintf("%s\nRepo: %s\nBuild: %s\nCommit: %s\nPipeline: %s", message, repoFullName, buildNumber, commit, bitbucketUrl)
}

func (slack *Slack) BuildBlocksByBitbucket(message string) *Slack {
	slack.Message = message

	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("%s\nCommit: %s", message, commit),
		},
	}

	deviderBlock := map[string]any{
		"type": "divider",
	}

	repoBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Repo: %s", repoFullName),
		},
	}

	goVersionBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Go version OS: %s", slack.GoVersion),
		},
	}

	goToolchainVersionBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Go toolchain version: %s", slack.GoToolchainVersion),
		},
	}

	slack.Blocks = append(slack.Blocks, messageBlock)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, repoBlock)
	slack.Blocks = append(slack.Blocks, goVersionBlock)
	slack.Blocks = append(slack.Blocks, goToolchainVersionBlock)

	dependencyUpdatesMsg := "no dependency updates needed"
	if len(slack.DependencyUpdates) > 0 {
		dependencyUpdatesMsg = "dependency updates needed: \n" + fmt.Sprintf("%s", slack.DependencyUpdates)
	}

	if len(slack.DependencyUpdates) > 40 {
		dependencyUpdatesMsg = "dependency updates needed: \n" + fmt.Sprintf("%s", slack.DependencyUpdates[:40]) + "\n...\n" +
			fmt.Sprintf("total of %d updates; have a look at the pipe", len(slack.DependencyUpdates))
	}

	dependencyUpdatesBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "Plain_text",
			"text": dependencyUpdatesMsg,
		},
	}
	slack.Blocks = append(slack.Blocks, dependencyUpdatesBlock)

	vulCheckBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": slack.VulCheck,
		},
	}

	slack.Blocks = append(slack.Blocks, vulCheckBlock)

	if origin == "" {
		return slack
	}

	actionBlock := map[string]any{
		"type": "actions",
		"elements": []map[string]any{
			{
				"type": "button",
				"text": map[string]any{
					"type": "plain_text",
					"text": fmt.Sprintf("Pipe %s", buildNumber),
				},
				"url": fmt.Sprintf("%s/addon/pipelines/home#!/results/%s", origin, commit),
			},
		},
	}
	slack.Blocks = append(slack.Blocks, actionBlock)

	return slack
}
