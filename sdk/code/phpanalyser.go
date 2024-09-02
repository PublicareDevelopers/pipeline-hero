package code

import "sync"

func NewPHPAnalyser() *PHPAnalyser {
	return &PHPAnalyser{
		lock: &sync.Mutex{},
	}
}

func (a *PHPAnalyser) SetThreshold(threshold float64) *PHPAnalyser {
	a.lock.Lock()
	a.Threshold = threshold
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) SetVulnCheck(vulnCheck string) *PHPAnalyser {
	a.lock.Lock()
	a.VulnCheck = vulnCheck
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) SetVulnCheckFail() *PHPAnalyser {
	a.lock.Lock()
	a.HasVulnCheckFail = true
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) SetCoverage(coverage float64) *PHPAnalyser {
	a.lock.Lock()
	a.Coverage = coverage
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) SetTestResult(testResult string) *PHPAnalyser {
	a.lock.Lock()
	a.TestResult = testResult
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) PushError(err string) *PHPAnalyser {
	a.lock.Lock()
	a.Errors = append(a.Errors, err)
	a.HasErrors = true
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) PushWarning(warn string) *PHPAnalyser {
	a.lock.Lock()
	a.Warnings = append(a.Warnings, warn)
	a.HasWarnings = true
	a.lock.Unlock()
	return a
}

func (a *PHPAnalyser) GetErrors() []string {
	return a.Errors
}

func (a *PHPAnalyser) GetCoverage() float64 {
	return a.Coverage
}

func (a *PHPAnalyser) GetTestResult() string {
	return a.TestResult
}

func (a *PHPAnalyser) GetVulnCheck() string {
	return a.VulnCheck
}
