package xio

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func QuietCloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		if !errors.Is(err, os.ErrClosed) {
			logrus.Errorf("encountered error while closing file %s: %s", file.Name(), err)
		}
	}
}
