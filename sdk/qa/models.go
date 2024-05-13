package qa

type ServerlessYML struct {
	FrameworkVersion string                                  `yaml:"frameworkVersion"`
	Provider         ServerlessProviderDefinition            `yaml:"provider"`
	Functions        map[string]ServerlessFunctionDefinition `yaml:"functions"`
	Resources        ServerlessResource                      `yaml:"resources"`
}

type ServerlessProviderDefinition struct {
	Name         string            `yaml:"name"`
	Stage        string            `yaml:"stage"`
	Runtime      string            `yaml:"runtime"`
	Architecture string            `yaml:"architecture"`
	Region       string            `yaml:"region"`
	Environment  map[string]string `yaml:"environment"`
}

type ServerlessFunctionDefinition struct {
	MemorySize  int                         `yaml:"memorySize"`
	Timeout     int                         `yaml:"timeout"`
	Handler     string                      `yaml:"handler"`
	Package     ServerlessPackageDefinition `yaml:"package"`
	Events      []interface{}               `yaml:"events"`
	Environment map[string]string           `yaml:"environment"`
}

type ServerlessResource struct {
	Definitions map[string]ServerlessResourceDefinition `yaml:"Resources"`
}

type ServerlessResourceDefinition struct {
	Type       string         `yaml:"Type"`
	Properties map[string]any `yaml:"Properties"`
}

type ServerlessPackageDefinition struct {
	Artifact string `yaml:"artifact"`
}

type ResourceInfo struct {
	Path string
	File string
}

type FunctionInfo struct {
	Path string
	File string
}
