package slack

import (
	"/pipeline-hero/sdk/cmds"
	"/pipeline-hero/sdk/code"
	"fmt"
	"os"
	"strings"
	"time"
)

var maxLengthDepUpdates = 20

func (client *Client) BuildThreadBlocks(analyser *code.Analyser) error {
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	client.Blocks = append(client.Blocks, getTestDurationBlock(analyser.GetProfiles()))
	client.Blocks = append(client.Blocks, getDividerBlock())

	toolchain := analyser.Module.Toolchain
	if toolchain != "" {
		goToolchainVersionBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "plain_text",
				"text": fmt.Sprintf("Go toolchain version: %s", toolchain),
			},
		}
		client.Blocks = append(client.Blocks, goToolchainVersionBlock)
		client.Blocks = append(client.Blocks, getDividerBlock())
	}

	if origin != "" {
		//make sure we have https, not only http
		if strings.HasPrefix(origin, "http://") {
			origin = strings.Replace(origin, "http://", "https://", 1)
		}

		pipeLink := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf("[Pipe #%s](%s)", buildNumber, fmt.Sprintf("%s/pipelines/results/%s", origin, buildNumber)),
			},
		}
		client.Blocks = append(client.Blocks, pipeLink)
	}

	dependencyUpdatesMsg := "no dependency updates needed"

	if len(analyser.Updates) > 0 {
		dependencyUpdatesMsg = "*dependency updates needed:* \n"
	}

	dependencyUpdatesBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": dependencyUpdatesMsg,
		},
	}
	client.Blocks = append(client.Blocks, dependencyUpdatesBlock)

	for _, updatableDependency := range analyser.Updates {
		dependencyUpdatesBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf("%s dependency update %s -> %s\n",
					updatableDependency.Path,
					updatableDependency.Version,
					updatableDependency.AvailableVersion),
			},
		}
		client.Blocks = append(client.Blocks, dependencyUpdatesBlock)
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
		client.Blocks = append(client.Blocks, warningsBlock)
	}

	if analyser.TestResult != "" {
		testResultBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Test result:*\n"),
			},
		}
		client.Blocks = append(client.Blocks, testResultBlock)
		client.Blocks = append(client.Blocks, getDividerBlock())

		//split analyser.TestResult in text blocks not longer than 3000 characters
		//slack has a limit of 3000 characters per text block

		testResultMsg := analyser.TestResult
		for len(testResultMsg) > 3000 {
			testResultBlock = map[string]any{
				"type": "section",
				"text": map[string]any{
					"type": "mrkdwn",
					"text": testResultMsg[:3000],
				},
			}
			client.Blocks = append(client.Blocks, testResultBlock)
			testResultMsg = testResultMsg[3000:]
		}

		testResultBlock = map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": testResultMsg,
			},
		}

		client.Blocks = append(client.Blocks, testResultBlock)
	}

	return nil
}

// BuildBlocks
// Deprecated: use BuildThreadBlocks instead
func (client *Client) BuildBlocks(analyser *code.Analyser) error {
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

	toolchain := analyser.Module.Toolchain
	if toolchain == "" {
		toolchain = "no toolchain found"
	}

	goToolchainVersionBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Go toolchain version: %s", toolchain),
		},
	}

	client.Blocks = append(client.Blocks, messageBlock)
	client.Blocks = append(client.Blocks, deviderBlock)
	client.Blocks = append(client.Blocks, getRepoBlock())
	client.Blocks = append(client.Blocks, getCommitMessageBlock())
	client.Blocks = append(client.Blocks, goVersionBlock)
	client.Blocks = append(client.Blocks, goToolchainVersionBlock)

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
	client.Blocks = append(client.Blocks, dependencyUpdatesBlock)

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
			client.Blocks = append(client.Blocks, vulnCheckBlock)
			vulnCheckMsg = vulnCheckMsg[3000:]
		}
		vulnCheckBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": vulnCheckMsg,
			},
		}
		client.Blocks = append(client.Blocks, vulnCheckBlock)
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
		client.Blocks = append(client.Blocks, warningsBlock)
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
		client.Blocks = append(client.Blocks, errorsBlock)
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
	client.Blocks = append(client.Blocks, pipeLink)

	return nil
}

func (client *Client) BuildJSBlocks(analyser *code.Analyser) error {
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	/*messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", analyser.GetCoverageInterpretation()),
		},
	}*/

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

	//client.Blocks = append(client.Blocks, messageBlock)
	//client.Blocks = append(client.Blocks, deviderBlock)
	client.Blocks = append(client.Blocks, npmMessageBlock)
	client.Blocks = append(client.Blocks, deviderBlock)
	client.Blocks = append(client.Blocks, getRepoBlock())
	client.Blocks = append(client.Blocks, deviderBlock)

	client.Blocks = append(client.Blocks, getCommitMessageBlock())

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
		client.Blocks = append(client.Blocks, warningsBlock)
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
			client.Blocks = append(client.Blocks, errorsBlock)
			msg = msg[3000:]
		}

		errorsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "plain_text",
				"text": msg,
			},
		}
		client.Blocks = append(client.Blocks, errorsBlock)
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
	client.Blocks = append(client.Blocks, pipeLink)

	return nil
}

func (client *Client) BuildPHPBlocks(analyser *code.Analyser) error {
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	commit := os.Getenv("BITBUCKET_COMMIT")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	/*messageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", analyser.GetCoverageInterpretation()),
		},
	}*/

	npmMessageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": "analyzed composer.json",
		},
	}

	deviderBlock := map[string]any{
		"type": "divider",
	}

	//client.Blocks = append(client.Blocks, messageBlock)
	//client.Blocks = append(client.Blocks, deviderBlock)
	client.Blocks = append(client.Blocks, npmMessageBlock)
	client.Blocks = append(client.Blocks, deviderBlock)
	client.Blocks = append(client.Blocks, getRepoBlock())
	client.Blocks = append(client.Blocks, deviderBlock)

	client.Blocks = append(client.Blocks, getCommitMessageBlock())

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
		client.Blocks = append(client.Blocks, warningsBlock)
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
			client.Blocks = append(client.Blocks, errorsBlock)
			msg = msg[3000:]
		}

		errorsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "plain_text",
				"text": msg,
			},
		}
		client.Blocks = append(client.Blocks, errorsBlock)
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
	client.Blocks = append(client.Blocks, pipeLink)

	return nil
}

func (client *Client) BuildErrorBlocks(analyser *code.Analyser, message string) error {
	buildNumber := os.Getenv("BITBUCKET_BUILD_NUMBER")
	origin := os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")

	deviderBlock := map[string]any{
		"type": "divider",
	}

	customMessageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", strings.Trim(message, "\n")),
		},
	}

	client.Blocks = append(client.Blocks, customMessageBlock)
	client.Blocks = append(client.Blocks, deviderBlock)
	client.Blocks = append(client.Blocks, getRepoBlock())
	client.Blocks = append(client.Blocks, getCommitMessageBlock())

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
		client.Blocks = append(client.Blocks, warningsBlock)
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
			client.Blocks = append(client.Blocks, errorsBlock)
			msg = msg[3000:]
		}

		errorsBlock := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "plain_text",
				"text": msg,
			},
		}
		client.Blocks = append(client.Blocks, errorsBlock)
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

	client.Blocks = append(client.Blocks, pipeLink)

	return nil
}

func (client *Client) GetBlocks() []map[string]any {
	return client.Blocks
}

func getRepoBlock() map[string]any {
	repoFullName := os.Getenv("BITBUCKET_REPO_FULL_NAME")
	branchName := os.Getenv("BITBUCKET_BRANCH")

	return map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("pipe run for repo: *%s*\nbranch: *%s*", strings.Trim(repoFullName, ""), strings.Trim(branchName, "")),
		},
	}
}

func getCommitMessageBlock() map[string]any {
	commitID := os.Getenv("BITBUCKET_COMMIT")
	commitMessage, err := cmds.GetCommitMessages()

	commitMessage = strings.Trim(commitMessage, "\n")

	if err != nil {
		return map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf("last commit %s \nerror getting commit message: %s", commitID, err.Error()),
			},
		}
	}

	return map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("last commit %s:\n%s", commitID, commitMessage),
		},
	}
}

func getDividerBlock() map[string]any {
	return map[string]any{
		"type": "divider",
	}
}

func getErrorsBlock(getErrors []string) ([]map[string]any, error) {
	returnBlocks := make([]map[string]any, 0)

	msg := "errors:\n"
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
		returnBlocks = append(returnBlocks, errorsBlock)
		msg = msg[3000:]
	}

	errorsBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": msg,
		},
	}
	returnBlocks = append(returnBlocks, errorsBlock)

	return returnBlocks, nil
}

func getTestDurationBlock(profiles []code.Profile) map[string]any {
	duration := time.Duration(0 * time.Second)
	for _, profile := range profiles {
		duration += time.Duration(profile.Duration * float64(time.Second))
	}

	return map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("test duration: %s", duration),
		},
	}

}
