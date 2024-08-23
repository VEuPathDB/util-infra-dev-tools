package compose

//
//     WARNING!!!
//
//
// For the types in this file to be more widely usable they must implement
// custom YAML deserialization methods as the compose spec allows for various
// config items to be of different types!
//
// An example case is the Service.Volumes property defined below.  A valid
// docker-compose file may have a list of values that could each individually
// be a string or a struct.
//

type Config struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services,omitempty"`
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
