package slack

type Slack struct {
	WebhookURL string
}

func New() *Slack {
	return &Slack{}
}
