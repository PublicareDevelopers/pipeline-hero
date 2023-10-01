package notifier

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"os"
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
}
