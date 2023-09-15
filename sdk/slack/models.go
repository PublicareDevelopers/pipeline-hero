package slack

type Slack struct {
	WebhookURL         string
	Message            string
	GoVersion          string
	GoToolchainVersion string
	DependencyUpdates  []string
	Errors             []string
	Blocks             []map[string]any
}

func New() *Slack {
	return &Slack{
		Blocks:            make([]map[string]any, 0),
		DependencyUpdates: make([]string, 0),
		Errors:            make([]string, 0),
	}
}
