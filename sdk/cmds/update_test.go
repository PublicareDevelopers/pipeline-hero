package cmds

import "testing"

func TestGetUpdateVersion(t *testing.T) {
	update, err := GetUpdateVersion("github.com/aws/aws-sdk-go@v1.45.11")
	if err != nil {
		t.Errorf("GetUpdateVersion() = %v, want %v", err, nil)
	}

	t.Log(update)
}
