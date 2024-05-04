package qa

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ZipPath struct {
	Path string
	File string
}

func CheckUnusedZips(rootDir string) (string, error) {
	var output strings.Builder

	// Step 1: Collect all zip paths from Makefiles
	zipPaths := make(map[string]ZipPath)
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".mk" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				zipPath, err := ExtractZipPath(line)
				if err == nil {
					zipPaths[zipPath] = ZipPath{Path: zipPath, File: path}
				}
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

	// Step 2: Check each zip path against all yml files
	err = filepath.Walk(filepath.Join(rootDir, "pipelines"), func(path string, info os.FileInfo, err error) error {
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
			for scanner.Scan() {
				line := scanner.Text()
				for zipPath := range zipPaths {
					// Create a regex to match the zip path
					re := regexp.MustCompile(regexp.QuoteMeta(zipPath))
					if re.MatchString(line) {
						delete(zipPaths, zipPath)
					}
				}
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

	// Step 3: Collect all unused zip paths
	for _, zipPath := range zipPaths {
		output.WriteString(fmt.Sprintf("Unused zip path: %s in %s\n", zipPath.Path, zipPath.File))
	}

	if output.Len() == 0 {
		return "No unused zip paths found", nil
	}

	return output.String(), fmt.Errorf("unused zip paths found")
}

// ExtractZipPath extracts the zip path from a makefile line.
func ExtractZipPath(line string) (string, error) {
	re := regexp.MustCompile(`zip -j\s+([^\s]+\.zip)`)
	match := re.FindStringSubmatch(line)
	if len(match) < 2 {
		return "", fmt.Errorf("no zip path found in line: %s", line)
	}
	return match[1], nil
}
