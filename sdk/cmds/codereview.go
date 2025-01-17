package cmds

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func CodeReview(codePart string) (string, error) {
	codeReview := ""
	errors := []error{}

	res, err := CodeReviewByStaticcheck(codePart)
	if err != nil {
		errors = append(errors, err)
	} else {
		codeReview += "\n"
		codeReview += res
	}

	res, err = CodeReviewByVet(codePart)
	if err != nil {
		errors = append(errors, err)
	} else {
		codeReview += "\n"
		codeReview += res
	}

	// TODO - this is not working as expected, we have no output from nilaway
	res, err = CodeReviewByNilCheck(codePart)
	if err != nil {
		errors = append(errors, err)
	} else {
		codeReview += "\n"
		type NilCheck map[string]any
		var nilCheck NilCheck

		err = json.Unmarshal([]byte(res), &nilCheck)
		if err != nil {
			errors = append(errors, err)
		}

		for _, check := range nilCheck {
			for k, v := range check.(map[string]any) {
				if k == "false" {
					continue
				}

				results, ok := (v.([]interface{}))
				if !ok {
					continue
				}

				for _, result := range results {
					nilaway, ok := result.(map[string]any)
					if !ok {
						continue
					}

					pos, ok := nilaway["posn"].(string)
					if !ok {
						continue
					}

					message, ok := nilaway["message"].(string)
					if !ok {
						continue
					}

					codeReview += fmt.Sprintf("at %s:\n %s\n\n", pos, message)
				}

			}

		}

		//codeReview += res
	}

	if len(errors) > 0 {
		return codeReview, fmt.Errorf("errors: %v", errors)
	}

	return codeReview, nil
}

func CodeReviewByStaticcheck(codePart string) (string, error) {
	_, err := exec.Command("go", "install", "honnef.co/go/tools/cmd/staticcheck@latest").Output()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("staticcheck", codePart).Output()
	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}

func CodeReviewByVet(codePart string) (string, error) {
	out, err := exec.Command("go", "vet", codePart).Output()
	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}

// CodeReviewByNilCheck runs nilaway on the given code part
func CodeReviewByNilCheck(codePart string) (string, error) {
	_, err := exec.Command("go", "install", "go.uber.org/nilaway/cmd/nilaway@latest").Output()
	if err != nil {
		return "", err
	}

	//-json -pretty-print=false
	out, err := exec.Command("nilaway", "-json", "-pretty-print", "false", codePart).Output()

	if err != nil {
		return string(out), nil
	}

	return string(out), nil
}

func CodeReviewVueMess(level string, output string, path string) (string, error) {
	_, err := exec.Command("npm", "install", "vue-mess-detector", "--save-dev").Output()
	if err != nil {
		return "", err
	}

	//npx vue-mess-detector analyze
	out, err := exec.Command("npx", "vue-mess-detector", "analyse",
		"--group=file", fmt.Sprintf("--level=%s", level),
		fmt.Sprintf("--path=%s", path),
		"--health-error=75", fmt.Sprintf("--output=%s", output),
	).Output()
	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}

func CodeReviewPHPStan(path string, level string) (string, error) {
	_, err := exec.Command("composer", "require", "phpstan/phpstan", "--dev").Output()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("vendor/bin/phpstan", "analyse", path, fmt.Sprintf("--level=%s", level)).Output()
	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}
