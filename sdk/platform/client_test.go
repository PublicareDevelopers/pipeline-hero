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

func TestClient_GetCommitAuthor(t *testing.T) {
	client := New()
	client.SetCommitAuthorRequest(CommitAuthorRequest{
		Repository: "vh-entity-interest",
		CommitId:   "b8f7c1f8b72cf82e9d3232f9c42ee2158ade4a83",
	},
	)

	author, err := client.GetCommitAuthor()
	if err != nil {
		t.Fatal(err)
	}

	if author != "Benjamin E." {
		t.Errorf("wanted to get author Benjamin E., got %s", author)
	}
}
