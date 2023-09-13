package slack

type Slack struct {
	WebhookURL string
	Blocks     []map[string]any
}

func New() *Slack {
	return &Slack{
		Blocks: make([]map[string]any, 0),
	}
}
