package slack

import (
	"errors"
	"os"
)

func (slack *Slack) InitConfiguration() (*Slack, error) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return slack, errors.New("SLACK_WEBHOOK_URL is not set")
	}

	slack.WebhookURL = webhookURL

	return slack, nil
}
