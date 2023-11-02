package cmds

import (
	"os"
	"os/exec"
)

func GetCommitMessages() (string, error) {
	commit := os.Getenv("BITBUCKET_COMMIT")
	if commit == "" {
		return "", nil
	}

	out, err := exec.Command("git", "log", "--format=%B", "-n", "1", commit).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
