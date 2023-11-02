package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"io/ioutil"
	"net/http"
	"os"
)

var maxLengthDepUpdates = 20

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
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", analyser.GetCoverageInterpretation()),
		},
	}

	deviderBlock := map[string]any{
		"type": "divider",
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
	slack.Blocks = append(slack.Blocks, getRepoBlock())
	slack.Blocks = append(slack.Blocks, getCommitMessageBlock())
	slack.Blocks = append(slack.Blocks, goVersionBlock)
	slack.Blocks = append(slack.Blocks, goToolchainVersionBlock)

	updatableDependencies := analyser.GetUpdatableDependencies()

	dependencyUpdatesMsg := "no dependency updates needed"
	if len(updatableDependencies) > 0 && len(updatableDependencies) <= maxLengthDepUpdates {
		dependencyUpdatesMsg = "dependency updates needed: \n"
		for _, updatableDependency := range updatableDependencies {
			dependencyUpdatesMsg +=
				fmt.Sprintf("* (used by %s) dependency update %s to %s\n",
					updatableDependency.From,
					updatableDependency.To,
					updatableDependency.UpdateTo)
		}
	}

	if len(updatableDependencies) > maxLengthDepUpdates {
		dependencyUpdatesMsg = "dependency updates needed: \n"
		for i, updatableDependency := range updatableDependencies {
			dependencyUpdatesMsg +=
				fmt.Sprintf("* (used by %s) dependency update %s to %s\n",
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
			"type": "mrkdwn",
			"text": dependencyUpdatesMsg,
		},
	}
	slack.Blocks = append(slack.Blocks, dependencyUpdatesBlock)

	//split analyser.VulnCheck in text blocks not longer than 3000 characters
	//slack has a limit of 3000 characters per text block
	if analyser.VulnCheck != "" {
		vulnCheckMsg := analyser.VulnCheck
		for len(vulnCheckMsg) > 3000 {
			vulnCheckBlock := map[string]any{
				"type": "section",
				"text": map[string]any{
					"type": "mrkdwn",
					"text": vulnCheckMsg[:3000],
				},
			}
			slack.Blocks = append(slack.Blocks, vulnCheckBlock)
			vulnCheckMsg = vulnCheckMsg[3000:]
		}
		vulnCheckBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": vulnCheckMsg,
			},
		}
		slack.Blocks = append(slack.Blocks, vulnCheckBlock)
	}

	warnings := analyser.GetWarnings()

	if len(warnings) > 0 {
		msg := "Warnings:\n"
		for _, warning := range warnings {
			msg += fmt.Sprintf(">%s\n", warning)
		}

		warningsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": msg,
			},
		}
		slack.Blocks = append(slack.Blocks, warningsBlock)
	}

	errors := analyser.GetErrors()

	if len(errors) > 0 {
		msg := "Errors:\n"
		for _, err := range errors {
			msg += fmt.Sprintf(">%s\n", err)
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

	pipeLink := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("[Pipe #%s](%s)", buildNumber, fmt.Sprintf("%s/addon/pipelines/home#!/results/%s", origin, commit)),
		},
	}
	slack.Blocks = append(slack.Blocks, pipeLink)

	return nil
}

func (slack *Slack) BuildJSBlocks(analyser *code.Analyser) error {
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", analyser.GetCoverageInterpretation()),
		},
	}

	npmMessageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": "analyzed package.json",
		},
	}

	deviderBlock := map[string]any{
		"type": "divider",
	}

	slack.Blocks = append(slack.Blocks, messageBlock)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, npmMessageBlock)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, getRepoBlock())
	slack.Blocks = append(slack.Blocks, deviderBlock)

	slack.Blocks = append(slack.Blocks, getCommitMessageBlock())

	warnings := analyser.GetWarnings()

	if len(warnings) > 0 {
		msg := "Warnings:\n"
		for _, warning := range warnings {
			msg += fmt.Sprintf(">%s\n", warning)
		}

		warningsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": msg,
			},
		}
		slack.Blocks = append(slack.Blocks, warningsBlock)
	}

	errors := analyser.GetErrors()

	if len(errors) > 0 {
		msg := "Errors:\n"
		for _, err := range errors {
			msg += fmt.Sprintf(">%s\n", err)
		}

		//split msg in text blocks not longer than 3000 characters
		//slack has a limit of 3000 characters per text block
		for len(msg) > 3000 {
			errorsBlock := map[string]any{
				"type": "section",
				"text": map[string]any{
					"type": "plain_text",
					"text": msg[:3000],
				},
			}
			slack.Blocks = append(slack.Blocks, errorsBlock)
			msg = msg[3000:]
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

	pipeLink := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("[Pipe #%s](%s)", buildNumber, fmt.Sprintf("%s/addon/pipelines/home#!/results/%s", origin, commit)),
		},
	}
	slack.Blocks = append(slack.Blocks, pipeLink)

	return nil
}

func (slack *Slack) BuildErrorBlocks(analyser *code.Analyser, message string) error {
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	deviderBlock := map[string]any{
		"type": "divider",
	}

	customMessageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", message),
		},
	}

	slack.Blocks = append(slack.Blocks, customMessageBlock)
	slack.Blocks = append(slack.Blocks, deviderBlock)
	slack.Blocks = append(slack.Blocks, getRepoBlock())
	slack.Blocks = append(slack.Blocks, getCommitMessageBlock())

	warnings := analyser.GetWarnings()

	if len(warnings) > 0 {
		msg := "Warnings:\n"
		for _, warning := range warnings {
			msg += fmt.Sprintf(">%s\n", warning)
		}

		warningsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": msg,
			},
		}
		slack.Blocks = append(slack.Blocks, warningsBlock)
	}

	getErrors := analyser.GetErrors()

	if len(getErrors) > 0 {
		msg := "Errors:\n"
		for _, err := range getErrors {
			msg += fmt.Sprintf(">%s\n", err)
		}

		//split msg in text blocks not longer than 3000 characters
		//slack has a limit of 3000 characters per text block
		for len(msg) > 3000 {
			errorsBlock := map[string]any{
				"type": "section",
				"text": map[string]any{
					"type": "plain_text",
					"text": msg[:3000],
				},
			}
			slack.Blocks = append(slack.Blocks, errorsBlock)
			msg = msg[3000:]
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

	commit := os.Getenv("BITBUCKET_COMMIT")

	pipeLink := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("[Pipe #%s](%s)", buildNumber, fmt.Sprintf("%s/addon/pipelines/home#!/results/%s", origin, commit)),
		},
	}

	slack.Blocks = append(slack.Blocks, pipeLink)

	return nil
}

func (slack *Slack) GetBlocks() []map[string]any {
	return slack.Blocks
}

func (slack *Slack) Notify() error {
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

func getRepoBlock() map[string]any {
	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	branchName := os.Getenv("BITBUCKET_BRANCH")

	return map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("Repo: *%s*\nBranch: %s", repoFullName, branchName),
		},
	}
}

func getCommitMessageBlock() map[string]any {
	commitID := os.Getenv("BITBUCKET_COMMIT")
	commitMessage, err := cmds.GetCommitMessages()
	if err != nil {
		return map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf("Commit %s \nError getting commit message: %s", commitID, err.Error()),
			},
		}
	}

	return map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("Commit %s:\n*%s*", commitID, commitMessage),
		},
	}
}
