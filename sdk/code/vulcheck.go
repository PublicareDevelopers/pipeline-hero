package code

import (
	"os/exec"
)

func VulCheck() (string, error) {
	_, err := exec.Command("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest").Output()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("govulncheck", "./...").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
