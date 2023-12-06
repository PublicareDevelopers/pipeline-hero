package platform

type Client struct {
	origin  string
	token   string
	request Request
}

type Request struct {
	Language string         `json:"language"`
	RunAt    string         `json:"runAt"`
	Build    int            `json:"build"`
	Analyser map[string]any `json:"analyser"`
	Context  Context        `json:"context"`
}

type Context struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	ThreadTs   string `json:"thread_ts"`
}

type Response struct {
}
