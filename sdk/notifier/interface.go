package notifier

import (
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
)

type Client interface {
	Validate() error
	BuildBlocks(analyser *code.Analyser) error
	BuildJSBlocks(analyser *code.Analyser) error
	GetBlocks() []map[string]any
	Notify() error
}
