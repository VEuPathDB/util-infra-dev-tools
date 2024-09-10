package use

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"vpdb-dev-tool/internal/lib/util"
)

type composeFile struct {
	Services map[string]composeService `yaml:"services"`
}

type composeService struct {
	Image string `yaml:"image"`
}

func parseModel(file *os.File) (config composeFile) {
	parser := yaml.NewDecoder(file)
	util.Must(parser.Decode(&config))
	logrus.Tracef("parsed docker compose file %s", file.Name())
	return
}
