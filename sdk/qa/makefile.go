package qa

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CheckLocalBuild(rootDir string) (string, error) {
	var output bytes.Buffer
	errorStatus := false

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".mk" {
			makefileDir := filepath.Dir(path)
			err := os.Chdir(makefileDir)
			if err != nil {
				return err
			}

			cmd := exec.Command("make", "-n", "-f", path)
			cmd.Stderr = &output
			err = cmd.Run()
			if err != nil {
				output.WriteString(fmt.Sprintf("Error in %s\n", path))
				errorStatus = true
			} else {
				cmd = exec.Command("make", "-f", path, "build")
				cmd.Stderr = &output
				err = cmd.Run()
				if err != nil {
					output.WriteString(fmt.Sprintf("Error in %s during build: %v\n", path, err))
					errorStatus = true
				}

				cmd = exec.Command("make", "-f", path, "zip")
				cmd.Stderr = &output
				err = cmd.Run()
				if err != nil {
					output.WriteString(fmt.Sprintf("Error in %s during zip: %v\n", path, err))
					errorStatus = true
				}

				if output.String() == "zip warning: name not matched" {
					output.WriteString(fmt.Sprintf("Warning in %s: name not matched\n", path))
					errorStatus = true
				}
			}

			err = os.Chdir(rootDir)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if errorStatus {
		return output.String(), fmt.Errorf("errors found in makefiles")
	}

	return output.String(), nil
}
