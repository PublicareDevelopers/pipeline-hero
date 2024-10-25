package npm

type Mod struct {
	Name         string                `json:"name"`
	Version      string                `json:"version"`
	Dependencies map[string]Dependency `json:"dependencies"`
}

type Dependency struct {
	Version    string `json:"version"`
	Resolved   string `json:"resolved"`
	Overridden bool   `json:"overridden"`
}
