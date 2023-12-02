package notifier

import (
	"encoding/json"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/slack"
	"os"
	"strings"
	"testing"
)

func TestSlack_Validate(t *testing.T) {
	handler, err := New("slack")
	if err != nil {
		t.Fatal(err)
	}

	err = handler.Client.Validate()
	if err == nil {
		t.Fatal("expected error")
	}

	_ = os.Setenv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX")
	err = handler.Client.Validate()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSlack_BuildBlocks(t *testing.T) {
	analyser := code.NewAnalyser().SetThreshold(50.0).
		SetGoVersion("go version go1.17.1 darwin/amd64").
		SetCoverageByTotal("total: (statements) 100.0%")

	handler, err := New("slack")
	if err != nil {
		t.Fatal(err)
	}

	err = handler.Client.BuildBlocks(analyser)
	if err != nil {
		t.Fatal(err)
	}

	blocks := handler.Client.GetBlocks()
	if len(blocks) == 0 {
		t.Fatalf("expected blocks, got 0")
	}

	blockStr, err := json.Marshal(blocks)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(blockStr), "go version go1.17.1 darwin/amd64") {
		t.Fatalf("expected go version, got %s", string(blockStr))
	}
}

func TestSlack_BuildErrorBlocks(t *testing.T) {
	analyser := code.NewAnalyser().SetThreshold(50.0).
		SetGoVersion("go version go1.17.1 darwin/amd64").
		SetCoverageByTotal("total: (statements) 100.0%").
		PushError("error occurred")

	handler, err := New("slack")
	if err != nil {
		t.Fatal(err)
	}

	err = handler.Client.BuildErrorBlocks(analyser, "pipe failed")
	if err != nil {
		t.Fatal(err)
	}

	blocks := handler.Client.GetBlocks()
	if len(blocks) == 0 {
		t.Fatalf("expected blocks, got 0")
	}

	blockStr, err := json.Marshal(blocks)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(blockStr), "pipe failed") {
		t.Fatalf("expected error, got %s", string(blockStr))
	}

	if !strings.Contains(string(blockStr), "error occurred") {
		t.Fatalf("expected error, got %s", string(blockStr))
	}
}

func TestSlack_SendErrorBlocks(t *testing.T) {
	analyser := code.NewAnalyser().SetThreshold(50.0).
		SetGoVersion("go version go1.20.13 darwin/amd64").
		SetCoverageByTotal("total: (statements) 100.0%").
		PushError("error occurred")

	handler, err := New("slack")
	if err != nil {
		t.Fatal(err)
	}

	err = handler.Client.BuildErrorBlocks(analyser, "pipe failed")
	if err != nil {
		t.Fatal(err)
	}

	blocks := handler.Client.GetBlocks()
	if len(blocks) == 0 {
		t.Fatalf("expected blocks, got 0")
	}

	slackClient, err := slack.NewTestClient()
	if err != nil {
		t.Fatal(err)
	}

	err = slackClient.SendProgressSlackBlocks(blocks)
	if err != nil {
		t.Fatal(err)
	}
}
