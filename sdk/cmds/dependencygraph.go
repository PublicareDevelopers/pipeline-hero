package cmds

import "os/exec"

func GetDependencyGraph() (string, error) {
	out, err := exec.Command("go", "mod", "graph").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
