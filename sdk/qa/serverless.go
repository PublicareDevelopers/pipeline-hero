package qa

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DefinitionInfo struct {
	Path     string
	File     string
	Errors   []string
	Warnings []string
	Def      ServerlessYML
}

type WarnLevels struct {
	FunctionsMemorySize int
}

var warnLevels = WarnLevels{
	FunctionsMemorySize: 4096,
}

func ServerlessQA(rootDir string) (map[string]DefinitionInfo, error) {
	// Schritt 1: Lesen Sie alle YML-Definitionen ein
	definitions := make(map[string]DefinitionInfo)
	err := ReadYMLDefinitions(rootDir, definitions)
	if err != nil {
		return definitions, err
	}

	// Schritt 2: Führen Sie die einzelnen Checks durch
	for path, definitionInfo := range definitions {
		// Check for runtime and architecture
		if definitionInfo.Def.Provider.Runtime != "provided.al2" || definitionInfo.Def.Provider.Architecture != "arm64" {
			definitionInfo.Warnings = append(definitionInfo.Warnings, "Warning: The combination of runtime: provided.al2 and architecture: arm64 is not found.")
		}

		for functionName, functionDefinition := range definitionInfo.Def.Functions {
			// Check for memory size of functions
			if functionDefinition.MemorySize > warnLevels.FunctionsMemorySize {
				definitionInfo.Warnings = append(definitionInfo.Warnings, "Warning: The memory size of function "+functionName+" is greater than 4096.")
			}

			//check if handler is bootstrap; when not
			// if runtime is provided.al2, handler should be bootstrap => add an error
			// else if handler is not bootstrap, but runtime is not provided.al2 => add a warning
			if definitionInfo.Def.Provider.Runtime == "provided.al2" && functionDefinition.Handler != "bootstrap" {
				definitionInfo.Errors = append(definitionInfo.Errors, "Error: The handler of function "+functionName+" should be bootstrap.")
			} else if definitionInfo.Def.Provider.Runtime != "provided.al2" && functionDefinition.Handler == "bootstrap" {
				definitionInfo.Warnings = append(definitionInfo.Warnings, "Warning: The handler of function "+functionName+" should not be bootstrap.")
			}

			//when events is empty => add a warning
			if len(functionDefinition.Events) == 0 {
				definitionInfo.Warnings = append(definitionInfo.Warnings, "Warning: The events of function "+functionName+" is empty.")
			}
		}

		definitions[path] = definitionInfo
	}

	return definitions, nil
}

func ReadYMLDefinitions(dir string, definitions map[string]DefinitionInfo) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

			//when there is not frameworkVersion in the yml file => no valid serverless yml file, we will skip
			if yml.FrameworkVersion == "" {
				return nil
			}

			definitions[path] = DefinitionInfo{Path: dir, File: path, Def: yml}
		}

		return nil
	})
}
