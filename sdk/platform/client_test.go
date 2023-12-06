package platform

import (
	"testing"
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
