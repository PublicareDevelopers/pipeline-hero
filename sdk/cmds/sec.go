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

func GoSASTProxy(exclude string, path string) (string, error) {
	_, err := exec.Command("go", "install", "github.com/securego/gosec/v2/cmd/gosec@latest").Output()
	if err != nil {
		return fmt.Sprintf("gosec not installed: %s", err), nil
	}

	var out []byte
	var outErr error

	if exclude == "" {
		fmt.Println("running gosec without exclude")
		out, outErr = exec.Command("gosec", "-enable-audit", "true", path).Output()
	} else {
		fmt.Println("running gosec with exclude", exclude)
		out, outErr = exec.Command("gosec", "-exclude", exclude, path).Output()

	}

	if outErr != nil {
		return "", errors.New(string(out))
	}

	return string(out), nil
}
