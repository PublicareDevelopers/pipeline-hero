package platform

type Client struct {
	origin              string
	token               string
	request             Request
	securityFixRequest  SecurityFixRequest
	sastFixRequest      SASTFixRequest
	dependenciesRequest DependenciesRequest
	commitAuthorRequest CommitAuthorRequest
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

type CommitAuthorRequest struct {
	BitbucketProject string `json:"bitbucketProject"`
	Repository       string `json:"repository"`
	CommitId         string `json:"commitId"`
}

type DependenciesRequest struct {
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Repository string         `json:"repository"`
	Name       string         `json:"name"`
	Version    string         `json:"version"`
	Language   string         `json:"language"`
	Data       map[string]any `json:"data"`
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
