package compose

type Config struct {
	Services map[string]Service `yaml:"services,omitempty"`
	Networks map[string]Network `yaml:"networks,omitempty"`
}

type Network struct {
	Aliases []string `yaml:"aliases,omitempty"`
}

type Service struct {
	// Build
	Image       string                `yaml:"image,omitempty"`
	DependsOn   map[string]Dependency `yaml:"depends_on,omitempty"`
	Entrypoint  string                `yaml:"entrypoint,omitempty"`
	HealthCheck *HealthCheck          `yaml:"healthcheck,omitempty"`
	Environment map[string]string     `yaml:"environment,omitempty"`
	Volumes     []Volume              `yaml:"volumes,omitempty"`
	Networks    map[string]Network    `yaml:"networks,omitempty"`
	Labels      map[string]string     `yaml:"labels,omitempty"`
}

type Volume struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

const DependencyConditionHealthy = "service_healthy"

type Dependency struct {
	Condition string `yaml:"condition"`
}

type HealthCheck struct {
	Test string `yaml:"test"`
}
