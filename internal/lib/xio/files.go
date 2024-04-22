package xio

import (
	"os"

	"github.com/sirupsen/logrus"
)

func MustStat(path string) os.FileInfo {
	stat, err := os.Stat(path)
	if err != nil {
		logrus.Fatalf("failed to stat file %s: %s", path, err)
	}
	return stat
}

func MustOpen(path string, flags int, perms os.FileMode) *os.File {
	handle, err := os.OpenFile(path, flags, perms)
	if err != nil {
		logrus.Fatalf("failed to open file %s: %s", path, err)
	}
	return handle
}
