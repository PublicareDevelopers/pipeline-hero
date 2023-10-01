package cmds

import "os/exec"

func GetVersion() (string, error) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
