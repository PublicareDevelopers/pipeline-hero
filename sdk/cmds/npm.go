package cmds

import (
	"os/exec"
)

func GetNPMAudit() (string, error) {
	out, err := exec.Command("npm", "audit").Output()
	if err != nil {
		return string(out), err
	}

	return string(out), err
}

func GetNPMOutdated() (string, error) {
	cmd := exec.Command("npm", "outdated", "--json")
	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return string(output), nil
			}
		}

		return string(output), err
	}

	return string(output), nil
}
