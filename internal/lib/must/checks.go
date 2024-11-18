package must

import (
	"github.com/sirupsen/logrus"
)

func NotError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

func Return1[T any](value T, err error) T {
	if err != nil {
		logrus.Fatal(err.Error())
	}

	return value
}
