package xio

import (
	"fmt"
	"io"
	"os"
)

type ReqRWFile struct {
	File *os.File
}

func (r ReqRWFile) Write(p []byte) (int, error) {
	if n, err := r.File.Write(p); err != nil {
		panic(fmt.Sprintf("encountered error while writing to file %s: %s\n", r.File.Name(), err))
	} else {
		return n, nil
	}
}

func (r ReqRWFile) WriteString(s string) (int, error) {
	if n, err := r.File.WriteString(s); err != nil {
		panic(fmt.Sprintf("encountered error while writing to file %s: %s\n", r.File.Name(), err))
	} else {
		return n, nil
	}
}

func (r ReqRWFile) Seek(offset int64, whence int) (int64, error) {
	if n, err := r.File.Seek(offset, whence); err != nil {
		if err == io.EOF {
			return n, err
		}

		panic(fmt.Sprintf("encountered error while seeking position in file %s: %s\n", r.File.Name(), err))
	} else {
		return n, nil
	}
}

func (r ReqRWFile) Read(p []byte) (int, error) {
	if n, err := r.File.Read(p); err != nil {
		if err == io.EOF {
			return n, err
		}

		panic(fmt.Sprintf("encountered error while reading file %s: %s\n", r.File.Name(), err))
	} else {
		return n, nil
	}
}
