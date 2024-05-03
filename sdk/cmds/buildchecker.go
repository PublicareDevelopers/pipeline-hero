package cmds

import (
	"embed"
	"io/ioutil"
	"os"
	"os/exec"
)

// make the shell directory as embedded

//go:embed bash/*
var bash embed.FS

func CheckBuild(rootDir string) (string, error) {
	yml, err := CheckYml(rootDir)
	if err != nil {
		return "", err
	}

	localBuild, err := CheckLocalBuild(rootDir)
	if err != nil {
		return "", err
	}

	return yml + "\n" + localBuild, nil
}

func CheckYml(rootDir string) (string, error) {
	f, err := bash.Open("bash/ymlcheck.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the file into a byte slice
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "ymlcheck.sh")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write the content to the temporary file
	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	// Run the script
	cmd := exec.Command("sh", tmpfile.Name(), rootDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func CheckLocalBuild(rootDir string) (string, error) {
	f, err := bash.Open("bash/localbuild.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the file into a byte slice
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "localbuild.sh")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write the content to the temporary file
	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	// Run the script
	cmd := exec.Command("sh", tmpfile.Name(), rootDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
