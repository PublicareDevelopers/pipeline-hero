package notifier

import "github.com/PublicareDevelopers/pipeline-hero/sdk/code"

type Slack struct{}

func (s *Slack) Validate() error {
	return nil
}

func (s *Slack) BuildBlocks(analyser *code.Analyser) error {
	return nil
}

func (s *Slack) Notify() error {
	return nil
}
