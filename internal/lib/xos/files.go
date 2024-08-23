package xos

import (
	"errors"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"vpdb-dev-tool/internal/lib/util"
)

func MustPathExists(path string) bool {
	return util.MustReturn(PathExists(path))
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func MustOpen(path string) *os.File {
	if file, err := os.Open(path); err != nil {
		logrus.Fatalf("failed to open file %s: %s", path, err)
		panic(err) // unreachable
	} else {
		logrus.Tracef("opened handle on file %s", path)
		return file
	}
}

func MustClose(file *os.File) {
	if err := file.Close(); err != nil {
		if !errors.Is(err, os.ErrClosed) {
			logrus.Errorf("failed to close handle on file %s: %s", file.Name(), err)
		}
	} else {
		logrus.Tracef("closed handle on file %s", file.Name())
	}
}

func MustCreateFile(name string) *os.File {
	if file, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644); err != nil {
		logrus.Fatalf("failed to create target file %s: %s", name, err)
		panic(err) // unreachable
	} else {
		return file
	}
}

func MustCopyFile(from, to string) {
	util.Must(CopyFile(from, to))
}

func CopyFile(from, to string) error {
	fstat, err := os.Stat(from)
	if err != nil {
		return err
	}

	f, err := os.Open(from)
	if err != nil {
		return err
	}
	defer MustClose(f)

	t, err := os.OpenFile(to, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fstat.Mode())
	if err != nil {
		return err
	}
	defer MustClose(t)

	_, err = io.Copy(t, f)

	return err
}
