package cmds

import (
	"errors"
	"fmt"
	"os/exec"
)

func VulnCheck(codePart string) (string, error) {
	_, err := exec.Command("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest").Output()
	if err != nil {
		return fmt.Sprintf("govulncheck not installed: %s", err), nil
	}

	out, err := exec.Command("govulncheck", codePart).Output()
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s", string(out)))
	}

	return string(out), nil
}
