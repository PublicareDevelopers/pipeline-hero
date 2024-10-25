package npm

type Mod struct {
	Name         string                `json:"name"`
	Version      string                `json:"version"`
	Dependencies map[string]Dependency `json:"dependencies"`
	Problems     []string              `json:"problems"`
}

type Dependency struct {
	Version    string `json:"version"`
	Required   string `json:"required"`
	Resolved   string `json:"resolved"`
	Overridden bool   `json:"overridden"`
}
