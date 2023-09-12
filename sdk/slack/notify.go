package slack

import (
	"bytes"
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

func BuildBitBucketMessage(message string) string {
	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")

	return fmt.Sprintf("%s\nRepo: %s\nBuild: %s\nCommit: %s", message, repoFullName, buildNumber, commit)
}
