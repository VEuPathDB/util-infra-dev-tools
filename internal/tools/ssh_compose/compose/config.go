package compose

type Config struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services,omitempty"`
	Volumes  []Volume           `yaml:"volumes,omitempty"`
	Networks map[string]Network `yaml:"networks,omitempty"`
}

type Network struct {
	Aliases []string `yaml:"aliases,omitempty"`
}

type Service struct {
	// Build
	Image       string             `yaml:"image,omitempty"`
	DependsOn   []string           `yaml:"depends_on,omitempty"`
	Entrypoint  []string           `yaml:"entrypoint,omitempty"`
	Environment map[string]string  `yaml:"environment,omitempty"`
	Volumes     []Volume           `yaml:"volumes,omitempty"`
	Networks    map[string]Network `yaml:"networks,omitempty"`
}

type Volume struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
