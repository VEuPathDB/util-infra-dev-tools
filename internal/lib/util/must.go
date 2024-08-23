package util

import "github.com/sirupsen/logrus"

func Must(err error) {
	MustMsg(err, "encountered unexpected error: %s")
}

func MustReturn[T any](value T, err error) T {
	Must(err)
	return value
}

func MustMsg(err error, fmt string) {
	if err != nil {
		logrus.Fatalf(fmt, err)
		panic(err) // unreachable
	}
}
