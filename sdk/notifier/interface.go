package notifier

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
)

type Client interface {
	Validate() error
	BuildBlocks(analyser *code.Analyser) error
	BuildJSBlocks(analyser *code.Analyser) error
	BuildPHPBlocks(analyser *code.Analyser) error
	BuildErrorBlocks(analyser *code.Analyser, message string) error
	GetBlocks() []map[string]any
	Notify() error
}
