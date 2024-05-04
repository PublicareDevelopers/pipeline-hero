package qa

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Task struct {
	Path string
}

func worker(tasks <-chan Task, wg *sync.WaitGroup, output *bytes.Buffer) {
	defer wg.Done()
	for task := range tasks {
		checkMakefile(task.Path, output)
	}
}

func checkMakefile(path string, output *bytes.Buffer) {
	makefileDir := filepath.Dir(path)
	err := os.Chdir(makefileDir)
	if err != nil {
		output.WriteString(fmt.Sprintf("Error in %s: %v\n", path, err))
		return
	}

	cmd := exec.Command("make", "-n", "-f", path)
	cmd.Stderr = output
	err = cmd.Run()
	if err != nil {
		output.WriteString(fmt.Sprintf("Error in %s\n", path))
		return
	}

	cmd = exec.Command("make", "-f", path, "build")
	cmd.Stderr = output
	err = cmd.Run()
	if err != nil {
		output.WriteString(fmt.Sprintf("Error in %s during build: %v\n", path, err))
		return
	}

	cmd = exec.Command("make", "-f", path, "zip")
	cmd.Stderr = output
	err = cmd.Run()
	if err != nil {
		output.WriteString(fmt.Sprintf("Error in %s during zip: %v\n", path, err))
		return
	}

	if output.String() == "zip warning: name not matched" {
		output.WriteString(fmt.Sprintf("Warning in %s: name not matched\n", path))
	}
}

func CheckLocalBuild(rootDir string) (string, error) {
	var output bytes.Buffer
	tasks := make(chan Task)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ { // Create 10 workers
		wg.Add(1)
		go worker(tasks, &wg, &output)
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".mk" {
			tasks <- Task{Path: path}
		}

		return nil
	})
	close(tasks)

	wg.Wait()

	if err != nil {
		return "", err
	}

	if output.Len() > 0 {
		return output.String(), fmt.Errorf("errors found in makefiles")
	}

	return output.String(), nil
}
