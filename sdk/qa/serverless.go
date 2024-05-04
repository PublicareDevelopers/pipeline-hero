package qa

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

type ServerlessCheck struct {
	Definitions map[string]DefinitionInfo
	MissingVars map[string][]string
}

func ServerlessQA(rootDir string, warnLevels WarnLevels) (ServerlessCheck, error) {
	serverlessCheck := ServerlessCheck{
		Definitions: make(map[string]DefinitionInfo),
		MissingVars: make(map[string][]string),
	}

	// Schritt 1: Lesen Sie alle YML-Definitionen ein
	definitions := make(map[string]DefinitionInfo)
	err := ReadYMLDefinitions(rootDir, definitions)
	if err != nil {
		return serverlessCheck, err
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
				definitionInfo.Warnings = append(definitionInfo.Warnings, fmt.Sprintf("Warning: The memory size of function %s is greater than %d.", functionName, warnLevels.FunctionsMemorySize))
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

	serverlessCheck.Definitions = definitions

	// Schritt 3: Überprüfen Sie die Codebestandteile
	err = AnalyseCode(rootDir, &serverlessCheck)
	if err != nil {
		return serverlessCheck, err
	}

	return serverlessCheck, nil
}

func AnalyseCode(rootDir string, check *ServerlessCheck) error {
	envVars := make(map[string][]string)

	// Step 1: Walk through all Go files
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".go" {
			// Step 2: Create the AST for the Go file
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}

			// Step 3: Walk through the AST and collect env vars
			ast.Inspect(f, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					if fun.Sel.Name == "Getenv" {
						if v, ok := call.Args[0].(*ast.BasicLit); ok {
							//remove the "" from the string
							envVar := v.Value[1 : len(v.Value)-1]

							//check if envVars[envVar] already exists
							if _, ok := envVars[envVar]; !ok {
								envVars[envVar] = []string{}
							}

							envVars[envVar] = append(envVars[envVar], path)
						}
					}
				}

				return true
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Step 4: Compare the collected env vars with the ones defined in the YML files
	missingVars := make(map[string][]string)
	for envVar, occurrences := range envVars {
		found := false
		for _, definitionInfo := range check.Definitions {
			for name := range definitionInfo.Def.Provider.Environment {
				if name == envVar {
					found = true
					break
				}
			}

			if found {
				break
			}

			for _, function := range definitionInfo.Def.Functions {
				for name := range function.Environment {
					if name == envVar {
						found = true
						break
					}
				}

				if found {
					break
				}
			}
		}

		if !found {
			missingVars[envVar] = occurrences
		}
	}

	check.MissingVars = missingVars

	return nil
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
