package sast

type SAST struct {
	GolangErrors struct {
	} `json:"Golang errors"`
	Issues []struct {
		Severity     string      `json:"severity"`
		Confidence   string      `json:"confidence"`
		Cwe          CWE         `json:"cwe"`
		RuleId       string      `json:"rule_id"`
		Details      string      `json:"details"`
		File         string      `json:"file"`
		Code         string      `json:"code"`
		Line         string      `json:"line"`
		Column       string      `json:"column"`
		Nosec        bool        `json:"nosec"`
		Suppressions interface{} `json:"suppressions"`
	} `json:"Issues"`
	Stats struct {
		Files int `json:"files"`
		Lines int `json:"lines"`
		Nosec int `json:"nosec"`
		Found int `json:"found"`
	} `json:"Stats"`
	GosecVersion string `json:"GosecVersion"`
}

type CWE struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func HasFailedSAST(sast SAST) (bool, error) {
	if sast.Stats.Found == 0 {
		return false, nil
	}

	return true, nil
}
