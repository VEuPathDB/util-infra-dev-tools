package use

import (
	"os"
	"vpdb-dev-tool/internal/lib/must"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type composeFile struct {
	Services map[string]composeService `yaml:"services"`
}

type composeService struct {
	Image string `yaml:"image"`
}

func parseModel(file *os.File) (config composeFile) {
	parser := yaml.NewDecoder(file)
	must.NotError(parser.Decode(&config))
	logrus.Tracef("parsed docker compose file %s", file.Name())
	return
}
