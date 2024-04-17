package platform

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()
	if client.origin == "" {
		t.Errorf("wanted to get origin from env")
	}

	if client.token == "" {
		t.Errorf("wanted to get token from env")
	}
}

func TestClient_Do(t *testing.T) {
	client := New()
	client.SetRequest(Request{
		Language: "go",
		RunAt:    time.Now().UTC().String(),
		Build:    1,
		Analyser: map[string]any{
			"unit test": "test",
		},
		Context: Context{
			Repository: "github.com/pipeline-hero-testings",
			Branch:     "main",
			ThreadTs:   "123456789",
		},
	})

	resp, err := client.Do()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}

func TestClient_CreateSecurityTask(t *testing.T) {
	client := New()
	client.SetSecurityFixRequest(SecurityFixRequest{
		Description:      "this is a security fix task from an unit test",
		BitbucketProject: "PHUTdeveloper",
		Repository:       "publicaremarketing/ph-unittest",
	})

	resp, err := client.CreateSecurityTask()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
