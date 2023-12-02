package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var webhookRatelimiter *rate.Limiter
var webhookRateLock *sync.Mutex

type Client struct {
	OAuthToken string `json:"oauth_token"`
	Channel    string `json:"channel"`
	ThreadTs   string `json:"thread_ts"`
	Errors     []string
	Blocks     []map[string]any
}

type WebhookMessage struct {
	Text string `json:"text"`
}

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type ResponseMessage struct {
	Channel  string `json:"channel"`
	ThreadTs string `json:"thread_ts"`
	Text     string `json:"text"`
}

type MessageResponse struct {
	Ok               bool                   `json:"ok"`
	Channel          string                 `json:"channel"`
	Ts               string                 `json:"ts"`
	Message          map[string]interface{} `json:"message"`
	Warning          string                 `json:"warning"`
	ResponseMetadata map[string]interface{} `json:"response_metadata"`
	Error            string                 `json:"error"`
}

func NewTestClient() (*Client, error) {
	oAuthToken := os.Getenv("SLACK_OAUTH_TOKEN")
	channel := os.Getenv("SLACK_MESSAGE_CHANNEL")
	if oAuthToken != "" || channel != "" {
		return &Client{
			OAuthToken: oAuthToken,
			Channel:    channel,
		}, nil
	}

	//check the actual directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//split after pipeline-hero/
	envDir := strings.SplitAfter(wd, "pipeline-hero/")[0]
	path := fmt.Sprintf("%s/.env", envDir)

	viper.SetConfigFile(path)

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading env file: %w", err)
	}

	return &Client{
		OAuthToken: viper.GetString("SLACK_OAUTH_TOKEN"),
		Channel:    viper.GetString("SLACK_MESSAGE_CHANNEL"),
	}, nil
}

func NewClient() (*Client, error) {
	oAuthToken := os.Getenv("SLACK_OAUTH_TOKEN")
	if oAuthToken == "" {
		return nil, fmt.Errorf("SLACK_OAUTH_TOKEN")
	}

	channel := os.Getenv("SLACK_MESSAGE_CHANNEL")
	if channel == "" {
		return nil, fmt.Errorf("SLACK_MESSAGE_CHANNEL")
	}

	return &Client{
		OAuthToken: oAuthToken,
		Channel:    channel,
	}, nil
}

func (client *Client) list() ([]any, error) {
	var bearer = "Bearer " + client.OAuthToken
	req, err := http.NewRequest("GET", "https://slack.com/api/conversations.list", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", bearer)

	var result map[string]any

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	channels, ok := result["channels"]
	if !ok {
		return nil, err
	}

	return channels.([]any), nil
}

func (client *Client) SendProgressSlackBlocks(blocks []map[string]any) error {
	err := client.sendSlackBlocks(client.Channel, blocks)
	if err != nil {
		log.Println("Error sending progress blocks to slack: ", err)
	}

	return err
}

func (client *Client) SendProgressSlackMessage(message string) error {
	err := client.sendSlackMessage(client.Channel, message)
	if err != nil {
		log.Println("Error sending progress message to slack: ", err)
		return err
	}

	return nil
}

func (client *Client) sendSlackBlocks(channelId string, blocks []map[string]any) error {
	message := map[string]any{
		"channel": channelId,
		"blocks":  blocks,
	}

	if client.ThreadTs != "" {
		message["thread_ts"] = client.ThreadTs
	}

	blockJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(blockJson)

	var bearer = "Bearer " + client.OAuthToken
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", reader)

	if err != nil {
		fmt.Printf("Error creating request: %s", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", bearer)

	httpClient := &http.Client{}

	ctx := context.Background()
	rl := getRateLimiter()
	fmt.Sprintf("Waiting for rate limiter")
	webhookRateLock.Lock()
	err = rl.Wait(ctx)
	webhookRateLock.Unlock()
	fmt.Sprintf("Now performing request")
	if err != nil {
		return err
	}

	response, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending message to Slack: %s", err)
		return err
	}

	var messageResponse MessageResponse
	err = json.NewDecoder(response.Body).Decode(&messageResponse)
	if err != nil {
		return err
	}

	if !messageResponse.Ok {
		return fmt.Errorf("error sending blocks to Slack: %s;\nmessage %s;\nresponse %+v", messageResponse.Warning, string(blockJson), messageResponse)
	}

	if client.ThreadTs == "" && messageResponse.Ts != "" {
		client.ThreadTs = messageResponse.Ts
	}
	return nil
}

func (client *Client) sendSlackMessage(channelId, message string) error {
	jsonValue := []byte(fmt.Sprintf(`{"channel": "%s", "text": "%s"}`, channelId, message))

	var bearer = "Bearer " + client.OAuthToken
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("Error creating request: %s", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", bearer)

	httpClient := &http.Client{}

	ctx := context.Background()
	rl := getRateLimiter()
	fmt.Sprintf("Waiting for rate limiter")
	webhookRateLock.Lock()
	err = rl.Wait(ctx)
	webhookRateLock.Unlock()
	fmt.Sprintf("Now performing request")
	if err != nil {
		return err
	}

	response, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending message to Slack: %s", err)
		return err
	}

	var messageResponse MessageResponse
	err = json.NewDecoder(response.Body).Decode(&messageResponse)
	if err != nil {
		return err
	}

	if !messageResponse.Ok {
		return fmt.Errorf("error sending message to Slack: %s;\nmessage %s;\nresponse %+v", messageResponse.Warning, string(jsonValue), messageResponse)
	}

	if client.ThreadTs == "" && messageResponse.Ts != "" {
		client.ThreadTs = messageResponse.Ts
	}
	return nil
}

func getRateLimiter() *rate.Limiter {
	if webhookRateLock == nil {
		webhookRateLock = &sync.Mutex{}
	}
	webhookRateLock.Lock()
	defer webhookRateLock.Unlock()

	if webhookRatelimiter == nil {
		webhookRatelimiter = rate.NewLimiter(rate.Every(time.Second), 1) // 1 request per second
	}
	return webhookRatelimiter
}
