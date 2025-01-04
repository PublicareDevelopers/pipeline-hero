package cmds

import (
	"testing"
)

func TestCodeReview(t *testing.T) {
	res, err := CodeReview("./...")
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	t.Log(res)
}
