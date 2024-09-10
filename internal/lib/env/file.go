package env

import (
	"os"

	"github.com/sirupsen/logrus"
)

const DotEnvFileName = ".env"

// GetOrCreateDotEnvFile returns a handle on a .env file in the current working
// directory, creating the file if it does not already exist.
//
// This function will panic if an error is encountered while attempting to open
// or create the file.
func GetOrCreateDotEnvFile() *os.File {
	if file, err := os.OpenFile(DotEnvFileName, os.O_CREATE, 0644); err != nil {
		logrus.Fatalf("failed to open "+DotEnvFileName+" file: %s", err)
		panic(err) // unreachable
	} else {
		return file
	}
}
