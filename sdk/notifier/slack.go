package notifier

import (
	"errors"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"os"
)

var maxLengthDepUpdates = 40

type Slack struct {
	WebhookURL string
	Errors     []string
	Blocks     []map[string]any
}

func (slack *Slack) Validate() error {
	slack.WebhookURL = os.Getenv("SLACK_WEBHOOK_URL")

	if slack.WebhookURL == "" {
		return errors.New("SLACK_WEBHOOK_URL is not set")
	}
	return nil
}

func (slack *Slack) BuildBlocks(analyser *code.Analyser) error {
	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("%s\nCommit: %s", analyser.GetCoverageInterpretation(), commit),
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
			"text": fmt.Sprintf("Go version OS: %s", analyser.GoVersion),
		},
	}

	goToolchainVersionBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Go toolchain version: %s", analyser.Toolchain),
		},
	}

	slack.Blocks = append(slack.Blocks, messageBlock)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, repoBlock)
	slack.Blocks = append(slack.Blocks, goVersionBlock)
	slack.Blocks = append(slack.Blocks, goToolchainVersionBlock)

	updatableDependencies := analyser.GetUpdatableDependencies()

	dependencyUpdatesMsg := "no dependency updates needed"
	if len(updatableDependencies) > 0 && len(updatableDependencies) <= maxLengthDepUpdates {
		dependencyUpdatesMsg = "dependency updates needed: \n"
		for _, updatableDependency := range updatableDependencies {
			dependencyUpdatesMsg +=
				fmt.Sprintf("(used by %s) dependency update %s to %s\n",
					updatableDependency.From,
					updatableDependency.To,
					updatableDependency.UpdateTo)
		}
	}

	if len(updatableDependencies) > maxLengthDepUpdates {

		for i, updatableDependency := range updatableDependencies {
			dependencyUpdatesMsg +=
				fmt.Sprintf("(used by %s) dependency update %s to %s\n",
					updatableDependency.From,
					updatableDependency.To,
					updatableDependency.UpdateTo)
			if i == maxLengthDepUpdates {
				break
			}
		}

		dependencyUpdatesMsg = fmt.Sprintf("total of %d updates; have a look at the pipe", len(updatableDependencies))
	}

	dependencyUpdatesBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": dependencyUpdatesMsg,
		},
	}
	slack.Blocks = append(slack.Blocks, dependencyUpdatesBlock)

	if analyser.VulnCheck != "" {
		vulCheckBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "plain_text",
				"text": analyser.VulnCheck,
			},
		}
		slack.Blocks = append(slack.Blocks, vulCheckBlock)
	}

	errors := analyser.GetErrors()

	if len(errors) > 0 {
		msg := "Errors:\n"
		for _, err := range errors {
			msg += fmt.Sprintf("%s\n", err)
		}

		errorsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "plain_text",
				"text": msg,
			},
		}
		slack.Blocks = append(slack.Blocks, errorsBlock)
	}

	if origin == "" {
		return nil
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

	return nil
}

func (slack *Slack) GetBlocks() []map[string]any {
	return slack.Blocks
}

func (slack *Slack) Notify() error {
	return nil
}
