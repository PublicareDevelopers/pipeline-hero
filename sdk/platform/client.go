package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func New() *Client {
	return &Client{
		token:  os.Getenv("PH_TOKEN"),
		origin: os.Getenv("PH_ORIGIN"),
	}
}

func (c *Client) SetRequest(request Request) {
	c.request = request
}

func (c *Client) SetSecurityFixRequest(request SecurityFixRequest) {
	c.securityFixRequest = request
}

func (c *Client) Do() (Response, error) {
	resp := Response{}

	payload, err := json.Marshal(c.request)
	if err != nil {
		return resp, err
	}

	fmt.Print("pushing data to platform")

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.origin+"push", bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating request: %s", err)
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s", err)
		return resp, err
	}

	if response.StatusCode >= 300 {
		return resp, fmt.Errorf("error status: %s", response.Status)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&resp)

	return resp, err
}

func (c *Client) CreateSecurityTask() (map[string]any, error) {
	resp := map[string]any{}

	payload, err := json.Marshal(c.securityFixRequest)
	if err != nil {
		return resp, err
	}

	fmt.Print("pushing data to platform")

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.origin+"security-fix", bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating request: %s", err)
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s", err)
		return resp, err
	}

	if response.StatusCode >= 300 {
		return resp, fmt.Errorf("error status: %s", response.Status)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&resp)

	return resp, err
}
