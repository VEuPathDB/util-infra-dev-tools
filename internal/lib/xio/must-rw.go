package xio

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type ReqRWFile struct {
	File *os.File
}

func (r ReqRWFile) WriteByte(b byte) error {
	buf := [1]byte{b}

	if _, err := r.File.Write(buf[:]); err != nil {
		logrus.Fatalf("encountered error while writing to file %s: %s", r.File.Name(), err)
		panic(nil) // unreachable
	}

	return nil
}

func (r ReqRWFile) Write(p []byte) (int, error) {
	if n, err := r.File.Write(p); err != nil {
		logrus.Fatalf("encountered error while writing to file %s: %s", r.File.Name(), err)
		panic(nil) // unreachable
	} else {
		return n, nil
	}
}

func (r ReqRWFile) WriteString(s string) (int, error) {
	if n, err := r.File.WriteString(s); err != nil {
		logrus.Fatalf("encountered error while writing to file %s: %s", r.File.Name(), err)
		panic(nil) // unreachable
	} else {
		return n, nil
	}
}

func (r ReqRWFile) WriteLine(s string) (int, error) {
	_, _ = r.WriteString(s)
	_ = r.WriteByte('\n')
	return 0, nil
}

func (r ReqRWFile) Seek(offset int64, whence int) (int64, error) {
	if n, err := r.File.Seek(offset, whence); err != nil {
		if err == io.EOF {
			return n, err
		}

		logrus.Fatalf("encountered error while seeking position in file %s: %s", r.File.Name(), err)
		panic(nil) // unreachable
	} else {
		return n, nil
	}
}

func (r ReqRWFile) Read(p []byte) (int, error) {
	if n, err := r.File.Read(p); err != nil {
		if err == io.EOF {
			return n, err
		}

		logrus.Fatalf("encountered error while reading file %s: %s", r.File.Name(), err)
		panic(nil) // unreachable
	} else {
		return n, nil
	}
}
