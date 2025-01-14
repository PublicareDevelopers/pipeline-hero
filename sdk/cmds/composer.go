package cmds

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func RunPHPUnitTest(phpunitCmd string, folder string) (string, error) {
	out, err := exec.Command(phpunitCmd, folder).Output()
	if err != nil {
		fmt.Printf("error: %s\n%s", err, string(out))
		return string(out), err
	}

	return string(out), err
}

func GetComposerOutDates() (string, error) {
	out, err := exec.Command("composer", "outdated", "-f", "json").Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return string(out), nil
			}
		}
		return string(out), err
	}

	return string(out), err
}

func GetComposerAudit() (string, error) {
	out, err := exec.Command("composer", "audit", "-f", "json").Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			//this is cause 6 means Abandoned only and at bitbucket 7 also
			if exitError.ExitCode() == 6 || exitError.ExitCode() == 7 {
				return parseComposerAudit(string(out)), nil
			}

			fmt.Println(exitError.ExitCode(), "code caused composer audit fail", err)
		}

		return parseComposerAudit(string(out)), err
	}

	return parseComposerAudit(string(out)), err
}

func parseComposerAudit(audit string) string {
	var msg string
	//avoid a panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			msg = audit
		}
	}()

	msg = parseJsonAudit(audit)

	return msg
}

func parseJsonAudit(audit string) string {
	var msg string
	type Audit struct {
		Advisories any `json:"advisories"` //can be [] or (map[string][]map[string]any
		Abandoned  any `json:"abandoned"`
	}

	var parsed Audit
	err := json.Unmarshal([]byte(audit), &parsed)
	if err != nil {
		fmt.Println("error:", err)
		return audit
	}
	//map[string][]map[string]any
	_, ok := parsed.Advisories.(map[string]any)
	if ok {
		type SafeAudit struct {
			Advisories map[string][]map[string]any `json:"advisories"` //can be [] or (map[string][]map[string]any
			Abandoned  map[string]any              `json:"abandoned"`
		}

		var safeAudit SafeAudit
		err := json.Unmarshal([]byte(audit), &safeAudit)
		if err != nil {
			fmt.Println("error:", err)
			return audit
		}

		for pck, advisories := range safeAudit.Advisories {
			msg += fmt.Sprintf("%s:\n", pck)
			for _, advisory := range advisories {
				msg += fmt.Sprintf("%s\nCVE: %s; Link: %s\nReportedAt: %s\nadvisoryId: %s; PackageName: %s; AffectedVersions: %s\n\n",
					advisory["title"], advisory["cve"], advisory["url"], advisory["reportedAt"], advisory["advisoryId"], advisory["packageName"], advisory["affectedVersions"])
				if advisory["sources"] == nil {
					continue
				}
				msg += "Sources:\n"
				for _, source := range advisory["sources"].([]interface{}) {
					source := source.(map[string]interface{})
					msg += fmt.Sprintf("%s:%s\n", source["name"], source["remoteId"])
				}
			}

			msg += "--------------------------------------------------\n"
		}
	}

	msg += "\n Abandoned \n"

	_, arrOk := parsed.Abandoned.([]any)
	if arrOk {
		return msg
	}

	parsedAbb, ok := parsed.Abandoned.(map[string]any)
	if !ok {
		return msg
	}

	for key, replace := range parsedAbb {
		msg += fmt.Sprintf("%s: %+v\n", key, replace)
	}

	return msg
}
