package code

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/npm"
	"sync"
)

func NewJSAnalyser() *JSAnalyser {
	return &JSAnalyser{
		lock: &sync.Mutex{},
	}
}

func (a *JSAnalyser) SetThreshold(threshold float64) *JSAnalyser {
	a.lock.Lock()
	a.Threshold = threshold
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) SetVulnCheck(vulnCheck string) *JSAnalyser {
	a.lock.Lock()
	a.VulnCheck = vulnCheck
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) SetCoverage(coverage float64) *JSAnalyser {
	a.lock.Lock()
	a.Coverage = coverage
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) SetTestResult(testResult string) *JSAnalyser {
	a.lock.Lock()
	a.TestResult = testResult
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) PushError(err string) *JSAnalyser {
	a.lock.Lock()
	a.Errors = append(a.Errors, err)
	a.HasErrors = true
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) GetErrors() []string {
	return a.Errors
}

func (a *JSAnalyser) GetCoverage() float64 {
	return a.Coverage
}

func (a *JSAnalyser) GetTestResult() string {
	return a.TestResult
}

func (a *JSAnalyser) GetVulnCheck() string {
	return a.VulnCheck
}

func (a *JSAnalyser) SetVulnCheckFail() *JSAnalyser {
	a.lock.Lock()
	a.HasVulnCheckFail = true
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) SetOutDates(outDates []npm.OutDate) *JSAnalyser {
	a.lock.Lock()
	a.OutDates = outDates
	a.lock.Unlock()
	return a
}

func (a *JSAnalyser) PushWarning(warning string) *JSAnalyser {
	a.lock.Lock()
	a.Warnings = append(a.Warnings, warning)
	a.HasWarnings = true
	a.lock.Unlock()
	return a
}
