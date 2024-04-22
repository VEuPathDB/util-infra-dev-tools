package xio

import (
	"errors"
	"io"
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

func QuietClose(io io.Closer) {
	if err := io.Close(); err != nil {
		if !errors.Is(err, os.ErrClosed) {
			logrus.Errorf("encountered error while closing unnamed stream: %s", err)
		}
	}
}
