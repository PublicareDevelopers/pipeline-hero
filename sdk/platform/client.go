package platform

import "os"

func New() *Client {
	return &Client{
		token:  os.Getenv("PH_TOKEN"),
		origin: os.Getenv("PH_ORIGIN"),
	}
}

func (c *Client) SetRequest(request Request) {
	c.request = request
}

func (c *Client) Do() (Response, error) {
	resp := Response{}

	return resp, nil
}
