package cmds

import (
	"fmt"
	"os/exec"
)

func Linter(codePart string) (string, error) {
	_, err := exec.Command("go", "install", "honnef.co/go/tools/cmd/staticcheck@latest").Output()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("staticcheck", codePart).Output()
	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}
