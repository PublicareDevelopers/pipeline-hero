package cmds

import (
	"os/exec"
	"regexp"
	"strings"
)

func GetUpdateVersion(original string) (string, error) {
	out, err := exec.Command("go", "list", "-m", "-u", original).Output()
	if err != nil {
		return "", err
	}

	res := string(out)

	//check if we have a construct like this: github.com/aws/aws-sdk-go v1.45.11 [v1.45.19] => we want to have the inside of []
	reg := regexp.MustCompile(`\[(.*)\]`)
	matches := reg.FindStringSubmatch(res)
	if len(matches) == 0 {
		return "", nil
	}

	//we have something like [v1.45.19]
	//we want to have v1.45.19
	update := matches[1]
	update = strings.Replace(update, "[", "", -1)
	update = strings.Replace(update, "]", "", -1)

	return update, nil
}
