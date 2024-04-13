package hosts

type Config struct {
	Hosts map[string][]string `yaml:"hosts"`
}
