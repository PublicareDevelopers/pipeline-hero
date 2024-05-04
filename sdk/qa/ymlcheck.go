package qa

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CheckYmlFiles(rootDir string) (string, error) {
	var output strings.Builder
	fileNotFound := false
	directories := []string{"pipelines/teststage", "pipelines/productivestage", "pipelines/infrastructure"}

	for _, dir := range directories {
		err := filepath.Walk(filepath.Join(rootDir, dir), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".yml" {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				lineNumber := 1
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, "artifact: bin/") {
						binPath := strings.TrimSpace(strings.Split(line, "artifact: bin/")[1])
						fullPath := filepath.Join(rootDir, "bin", binPath)
						if _, err := os.Stat(fullPath); os.IsNotExist(err) {
							output.WriteString(fmt.Sprintf("File does not exist: %s at line %d in %s\n", fullPath, lineNumber, path))
							fileNotFound = true
						}
					}
					lineNumber++
				}

				if err := scanner.Err(); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return "", err
		}
	}

	if fileNotFound {
		return output.String(), fmt.Errorf("files not found in yml files")
	}

	return output.String(), nil
}
