package xos

import (
	"errors"
	"io"
	"os"
	"vpdb-dev-tool/internal/lib/must"

	"github.com/sirupsen/logrus"
)

// MustPathExists wraps PathExists and panics on error.
func MustPathExists(path string) bool {
	return must.Return1(PathExists(path))
}

// PathExists tests whether a given path exists on the system.
func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func MustOpen(path string, mode int, perms os.FileMode) *os.File {
	return must.Return1(os.OpenFile(path, mode, perms))
}

// MustOpenSimple attempts to open a target file and panics on failure.
func MustOpenSimple(path string) *os.File {
	if file, err := os.Open(path); err != nil {
		logrus.Fatalf("failed to open file %s: %s", path, err)
		panic(err) // unreachable
	} else {
		logrus.Tracef("opened handle on file %s", path)
		return file
	}
}

// MustClose attempts to close the given file and logs an error on failure.
func MustClose(file io.Closer) {
	if err := file.Close(); err != nil {
		if !errors.Is(err, os.ErrClosed) {
			if v, ok := file.(*os.File); ok {
				logrus.Errorf("failed to close handle on file %s: %s", v.Name(), err)
			} else {
				logrus.Errorf("failed to close stream: %s", err.Error())
			}
		}
	} else {
		if v, ok := file.(*os.File); ok {
			logrus.Tracef("closed handle on file %s", v.Name())
		} else {
			logrus.Trace("closed handle on stream")
		}
	}
}

// MustCreateFile attempts to create a new file and panics if the file already
// exists or if creating a new file failed.
func MustCreateFile(name string) *os.File {
	if file, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644); err != nil {
		logrus.Fatalf("failed to create target file %s: %s", name, err)
		panic(err) // unreachable
	} else {
		return file
	}
}

// MustCopyFile wraps CopyFile and panics on error.
func MustCopyFile(from, to string) {
	must.NotError(CopyFile(from, to))
}

// CopyFile copies a file from one path to another.
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

func MustDelete(path string) {
	must.NotError(os.Remove(path))
}
