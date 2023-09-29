package cmds

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version, err := GetVersion()
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	if version == "" {
		t.Errorf("Error: version is empty\n")
	}

	if strings.Contains(version, "go version") == false {
		t.Errorf("Error: version is not correct\n")
	}
}
