package platform

type Client struct {
	origin             string
	token              string
	request            Request
	securityFixRequest SecurityFixRequest
	sastFixRequest     SASTFixRequest
}

type SecurityFixRequest struct {
	Description      string `json:"description"`
	BitbucketProject string `json:"bitbucketProject"`
	Repository       string `json:"repository"`
	Language         string `json:"language"`
}

type SASTFixRequest struct {
	Description      string `json:"description"`
	BitbucketProject string `json:"bitbucketProject"`
	Repository       string `json:"repository"`
	Language         string `json:"language"`
}

type Request struct {
	Language string         `json:"language"`
	RunAt    string         `json:"runAt"`
	Build    int            `json:"build"`
	Analyser map[string]any `json:"analyser"`
	Context  Context        `json:"context"`
}

type Context struct {
	Project    string `json:"project"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	ThreadTs   string `json:"thread_ts"`
}

type Response struct {
}
