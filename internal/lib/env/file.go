package env

import (
	"os"

	"github.com/sirupsen/logrus"
)

const DotEnvFileName = ".env"

func GetOrCreateDotEnvFile() *os.File {
	if file, err := os.OpenFile(DotEnvFileName, os.O_CREATE, 0644); err != nil {
		logrus.Fatalf("failed to open "+DotEnvFileName+" file: %s", err)
		panic(err) // unreachable
	} else {
		return file
	}
}
