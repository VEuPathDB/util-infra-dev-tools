package dh

import (
	"github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/sirupsen/logrus"
)

func TestTag(image, tag string) (bool, error) {
	logrus.Debugf("testing dockerhub for the existence of image %s/%s:%s", namespace, image, tag)

	res := simple.HeadRequest(makeTagURL(image, tag)).Submit()
	defer res.Close()

	if code, err := res.GetResponseCode(); err != nil {
		logrus.Errorf("request to dockerhub api failed with error: %s", err)
		return false, err
	} else {
		found := code == 200
		logrus.Debugf("docker image %s/%s:%s exists: %t", namespace, image, tag, found)
		return found, nil
	}
}
