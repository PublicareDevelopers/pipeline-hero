package slack

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
)

func (client *Client) StartConversation(analyser *code.Analyser, pipeType string) error {
	startBlocks := make([]map[string]any, 0)
	pipeErrors := analyser.GetErrors()

	repoBlock := getRepoBlock()

	startBlocks = append(startBlocks, repoBlock)
	startBlocks = append(startBlocks, getDividerBlock())

	if len(pipeErrors) > 0 {
		errorMessage := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf(":fire: *%s* \nactions required", "pipeline-hero failed"),
			},
		}

		startBlocks = append(startBlocks, errorMessage)
		startBlocks = append(startBlocks, getDividerBlock())
	} else {
		successMessage := map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf(":tadta: *%s* \nno urgent action required", "pipeline-hero success"),
			},
		}

		startBlocks = append(startBlocks, successMessage)
		startBlocks = append(startBlocks, getDividerBlock())
	}

	startBlocks = append(startBlocks, getCommitMessageBlock())
	startBlocks = append(startBlocks, getDividerBlock())

	goVersionBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "plain_text",
			"text": fmt.Sprintf("Go version OS: %s", analyser.GoVersion),
		},
	}

	startBlocks = append(startBlocks, goVersionBlock)
	startBlocks = append(startBlocks, getDividerBlock())

	coverageBlock := map[string]any{
		"type": "section",
		"text": map[string]any{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*", analyser.GetCoverageInterpretation()),
		},
	}

	startBlocks = append(startBlocks, coverageBlock)
	startBlocks = append(startBlocks, getDividerBlock())

	if len(pipeErrors) > 0 {
		errorBlocks, err := getErrorsBlock(pipeErrors)
		if err != nil {
			return err
		}

		startBlocks = append(startBlocks, errorBlocks...)
	}

	err := client.SendProgressSlackBlocks(startBlocks)
	if err != nil {
		return err
	}

	//now we can use the threads
	err = client.BuildThreadBlocks(analyser)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = client.SendProgressSlackBlocks(client.Blocks)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
