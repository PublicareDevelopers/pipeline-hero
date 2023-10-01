package code

import "fmt"

func (a *Analyser) CheckThreshold() error {
	if a.Coverage >= a.Threshold {
		return nil
	}

	return fmt.Errorf("coverage %.2f%% is below threshold %.2f%%", a.Coverage, a.Threshold)
}
