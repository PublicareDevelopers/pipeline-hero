package cmds

import "os/exec"

func SelfUpdate() (string, error) {
	out, err := exec.Command("go", "install", "github.com/PublicareDevelopers/pipeline-hero@latest").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
