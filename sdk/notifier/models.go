package notifier

import "fmt"

type Handler struct {
	Client Client
}

func New(client string) (*Handler, error) {
	switch client {
	case "slack":
		return &Handler{
			Client: &Slack{},
		}, nil
	default:
		return nil, fmt.Errorf("client %s not found", client)
	}
}
