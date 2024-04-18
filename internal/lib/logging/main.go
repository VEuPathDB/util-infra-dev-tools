package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogging(lvl logrus.Level) {
	logrus.SetLevel(lvl)
	logrus.SetOutput(os.Stderr)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
}
