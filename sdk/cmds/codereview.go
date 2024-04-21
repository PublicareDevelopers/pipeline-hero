package cmds

import (
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

	res, err = CodeReviewByGoSec(codePart)
	if err != nil {
		errors = append(errors, err)
	} else {
		codeReview += "\n"
		codeReview += res
	}

	// TODO - this is not working as expected, we have no output from nilaway
	//res, err = CodeReviewByNilCheck(codePart)
	//if err != nil {
	//	errors = append(errors, err)
	//} else {
	//	codeReview += "\n"
	//	codeReview += res
	//}

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

func CodeReviewByGoSec(codePart string) (string, error) {
	_, err := exec.Command("go", "install", "github.com/securego/gosec/v2/cmd/gosec@latest").Output()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("gosec", "-fmt", "txt", codePart).Output()
	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}

// CodeReviewByNilCheck runs nilaway on the given code part
// TODO - this is not working as expected, we have no output from nilaway
func CodeReviewByNilCheck(codePart string) (string, error) {
	_, err := exec.Command("go", "install", "go.uber.org/nilaway/cmd/nilaway@latest").Output()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("nilaway", codePart).Output()

	if err != nil {
		return fmt.Sprintf("%s", string(out)), nil
	}

	return string(out), nil
}
