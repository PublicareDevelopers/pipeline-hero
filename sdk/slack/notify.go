package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	blockJson, err := json.Marshal(slack.Blocks)
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
		return errors.New(fmt.Sprintf("Slack response status code: %d", resp.StatusCode))
	}

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
	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("%s\nRepo: Commit: %s", message, commit),
		},
	}

	deviderBlock := map[string]any{
		"type": "divider",
	}

	actionBlock := map[string]any{
		"type": "actions",
		"elements": []any{
			map[string]any{
				"type": "button",
				"text": fmt.Sprintf("Pipe %s", buildNumber),
				"url":  fmt.Sprintf("%s/addon/pipelines/home#!/results/%s", origin, commit),
			},
		},
	}

	headerBlcok := map[string]any{
		"type": "header",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Repo: %s", repoFullName),
		},
	}

	slack.Blocks = append(slack.Blocks, headerBlcok)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, messageBlock)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, actionBlock)

	return slack
}