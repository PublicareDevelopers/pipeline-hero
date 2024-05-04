package cmds

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/qa"
)

func CheckBuild(rootDir string) (string, error) {
	yml, err := CheckYml(rootDir)
	if err != nil {
		return "", err
	}

	localBuild, err := CheckLocalBuild(rootDir)
	if err != nil {
		return "", err
	}

	unusedZips, err := CheckUnusedZips(rootDir)
	if err != nil {
		return "", err
	}

	return yml + "\n" +
		localBuild + "\n" +
		unusedZips, nil
}

func CheckYml(rootDir string) (string, error) {
	return qa.CheckYmlFiles(rootDir)
}

func CheckLocalBuild(rootDir string) (string, error) {
	return qa.CheckLocalBuild(rootDir)
}

func CheckUnusedZips(rootDir string) (string, error) {
	return qa.CheckUnusedZips(rootDir)
}

func CheckFunctionDefinitions(rootDir string) (string, error) {
	return qa.CheckFunctionDefinitions(rootDir)
}
