package qa

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CheckFunctionDefinitions(rootDir string) (string, error) {
	var output strings.Builder

	// Step 1: Collect all function names from teststage yml files
	testFunctions := make(map[string]FunctionInfo)
	err := filepath.Walk(filepath.Join(rootDir, "pipelines", "teststage"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yml" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var yml ServerlessYML
			err = yaml.Unmarshal(file, &yml)
			if err != nil {
				return err
			}

			for functionName := range yml.Functions {
				testFunctions[functionName] = FunctionInfo{Path: functionName, File: path}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Step 2: Check each function name against all productivestage yml files
	err = filepath.Walk(filepath.Join(rootDir, "pipelines", "productivestage"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yml" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var yml ServerlessYML
			err = yaml.Unmarshal(file, &yml)
			if err != nil {
				return err
			}

			for functionName := range yml.Functions {
				if _, ok := testFunctions[functionName]; ok {
					delete(testFunctions, functionName)
				}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Step 3: Collect all missing function names
	for _, functionInfo := range testFunctions {
		output.WriteString(fmt.Sprintf("Missing productive function: %s, defined in %s\n", functionInfo.Path, functionInfo.File))
	}

	if output.Len() == 0 {
		return "No missing functions found", nil
	}

	return output.String(), fmt.Errorf("missing functions found")
}

func CheckResourceDefinitions(rootDir string) (string, error) {
	var output strings.Builder

	// Step 1: Collect all resource names from teststage yml files
	testResources := make(map[string]ResourceInfo)
	err := filepath.Walk(filepath.Join(rootDir, "pipelines", "teststage"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yml" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var yml ServerlessYML
			err = yaml.Unmarshal(file, &yml)
			if err != nil {
				return err
			}

			for resourceName := range yml.Resources.Definitions {
				testResources[resourceName] = ResourceInfo{Path: resourceName, File: path}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Step 2: Check each resource name against all productivestage yml files
	err = filepath.Walk(filepath.Join(rootDir, "pipelines", "productivestage"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yml" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var yml ServerlessYML
			err = yaml.Unmarshal(file, &yml)
			if err != nil {
				return err
			}

			for resourceName := range yml.Resources.Definitions {
				if _, ok := testResources[resourceName]; ok {
					delete(testResources, resourceName)
				}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	output.WriteString(fmt.Sprintf("checking %d resources\n", len(testResources)))

	// Step 3: Collect all missing resource names
	for _, resourceInfo := range testResources {
		output.WriteString(fmt.Sprintf("Missing productive resource: %s, defined in %s\n", resourceInfo.Path, resourceInfo.File))
	}

	if output.Len() == 0 {
		return "No missing resources found", nil
	}

	return output.String(), fmt.Errorf("missing resources found")
}
