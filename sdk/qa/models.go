package qa

type ServerlessYML struct {
	FrameworkVersion string                                  `yaml:"frameworkVersion"`
	Provider         ServerlessProviderDefinition            `yaml:"provider"`
	Functions        map[string]ServerlessFunctionDefinition `yaml:"functions"`
}

type ServerlessProviderDefinition struct {
	Name         string `yaml:"name"`
	Stage        string `yaml:"stage"`
	Runtime      string `yaml:"runtime"`
	Architecture string `yaml:"architecture"`
	Region       string `yaml:"region"`
	Environment  map[string]string
}

type ServerlessFunctionDefinition struct {
	MemorySize int                         `yaml:"memorySize"`
	Timeout    int                         `yaml:"timeout"`
	Handler    string                      `yaml:"handler"`
	Package    ServerlessPackageDefinition `yaml:"package"`
	Events     []interface{}               `yaml:"events"`
}

type ServerlessPackageDefinition struct {
	Artifact string `yaml:"artifact"`
}

type FunctionInfo struct {
	Path string
	File string
}
