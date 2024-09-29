package cmds

import (
	"errors"
	"fmt"
	"os/exec"
)

func GoSAST() (string, error) {
	_, err := exec.Command("go", "install", "github.com/securego/gosec/v2/cmd/gosec@latest").Output()
	if err != nil {
		return fmt.Sprintf("gosec not installed: %s", err), nil
	}

	out, err := exec.Command("gosec", "-fmt", "json", "-enable-audit", "true", "./...").Output()
	if err != nil {
		return "", errors.New(string(out))
	}

	return string(out), nil
}
