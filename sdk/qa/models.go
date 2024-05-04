package qa

type ServerlessYML struct {
	Functions map[string]interface{} `yaml:"functions"`
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
