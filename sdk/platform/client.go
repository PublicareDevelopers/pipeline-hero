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

func (c *Client) SetSASTFixRequest(request SASTFixRequest) {
	c.sastFixRequest = request
}

func (c *Client) SetDependenciesRequest(request DependenciesRequest) {
	c.dependenciesRequest = request
}

func (c *Client) SetCommitAuthorRequest(request CommitAuthorRequest) {
	c.commitAuthorRequest = request
}

func (c *Client) Do() (Response, error) {
	resp := Response{}

	payload, err := json.Marshal(c.request)
	if err != nil {
		return resp, err
	}

	fmt.Println("pushing analyser data to platform")

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

	fmt.Println("check to create a sec task by platform")

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

func (c *Client) CreateSASTTask() (map[string]any, error) {
	resp := map[string]any{}

	payload, err := json.Marshal(c.sastFixRequest)
	if err != nil {
		return resp, err
	}

	fmt.Println("check to create a SAST task by platform")

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.origin+"sast-fix", bytes.NewReader(payload))
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

func (c *Client) SendDependencies() (map[string]any, error) {
	resp := map[string]any{}

	payload, err := json.Marshal(c.dependenciesRequest)
	if err != nil {
		return resp, err
	}

	fmt.Println("pushing dependencies data to platform")

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.origin+"dependencies", bytes.NewReader(payload))
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

func (c *Client) GetCommitAuthor() (string, error) {
	author := "channel"

	payload, err := json.Marshal(c.commitAuthorRequest)
	if err != nil {
		return author, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.origin+"commitauthor", bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating request: %s", err)
		return author, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s", err)
		return author, err
	}

	if response.StatusCode >= 300 {
		var errResp any
		err = json.NewDecoder(response.Body).Decode(&errResp)

		return author, fmt.Errorf("error status: %s; %+v", response.Status, errResp)
	}

	defer response.Body.Close()
	type AuthorResponse struct {
		Author string `json:"author"`
	}

	var authResp AuthorResponse

	err = json.NewDecoder(response.Body).Decode(&authResp)
	if err != nil {
		return author, err
	}

	return authResp.Author, nil
}
