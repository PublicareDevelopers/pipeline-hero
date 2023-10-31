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
