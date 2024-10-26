package cmds

import (
	"fmt"
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

func AnalyseNPMModule() (string, error) {
	out, err := exec.Command("npm", "list", "--json").Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func GetNPMPackageRepoURL(packageName string) (string, error) {
	out, err := exec.Command("npm", "view", packageName, "repository.url").Output()
	if err != nil {
		fmt.Println("Error getting npm package repository URL: ", err)
		return string(out), err
	}

	return string(out), nil
}

func GetNPMPackageContributors(packageName string) (string, error) {
	out, err := exec.Command("npm", "view", packageName, "contributors").Output()
	if err != nil {
		fmt.Println("Error getting npm package contributors: ", err)
		return string(out), err
	}

	return string(out), nil
}

func GetNPMAuthor(packageName string) (string, error) {
	out, err := exec.Command("npm", "view", packageName, "author").Output()
	if err != nil {
		fmt.Println("Error getting npm package author: ", err)
		return string(out), err
	}

	return string(out), nil
}

func GetNPMPackageLicense(packageName string) (string, error) {
	out, err := exec.Command("npm", "view", packageName, "license").Output()
	if err != nil {
		fmt.Println("Error getting npm package license: ", err)
		return string(out), err
	}

	return string(out), nil
}
