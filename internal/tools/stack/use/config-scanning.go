package use

import (
	"strings"

	"github.com/sirupsen/logrus"

	"vpdb-dev-tool/internal/lib/env"
	"vpdb-dev-tool/internal/lib/xos"
)

func scanConfigs(files []string) map[string]string {
	images := make(map[string]string, 8)

	for _, file := range files {
		logrus.Debugf("scanning file %s", file)
		scanConfig(file, parseConfig(file), images)
	}

	return images
}

func scanConfig(file string, config composeFile, index map[string]string) {
	var splitPos int
	var imageName, tagName string

	logrus.Debugf("file %s contained %d service entries", file, len(config.Services))

	for name, service := range config.Services {
		if len(service.Image) == 0 {
			logrus.Debugf("%s service %s specifies no image, skipping", file, name)
			continue
		}

		if splitPos = strings.IndexByte(service.Image, ':'); splitPos < 1 {
			logrus.Debugf("%s service %s image doesn't appear to have a tag specification, skipping", file, name)
			continue
		}

		if len(service.Image) == splitPos+1 {
			logrus.Warningf("%s service %s image ends in ':', probably invalid", file, name)
			continue
		}

		if service.Image[splitPos+1] != '$' {
			logrus.Debugf("%s service %s image specifies an unmodifiable tag", file, name)
			continue
		}

		imageName = trimImageName(service.Image[:splitPos])

		if service.Image[splitPos+2] != '{' {
			logrus.Warningf("%s service %s image has wonky tag, skipping", file, name)
			continue
		}

		tagName = service.Image[splitPos+3:]

		if splitPos = strings.IndexAny(tagName, ":?}"); splitPos < 1 {
			logrus.Warningf("%s service %s image has malformed tag variable", file, name)
			continue
		}

		tagName = tagName[:splitPos]

		if !env.ResemblesEnvKey(tagName) {
			logrus.Warningf("%s service %s image has malformed tag variable", file, name)
			continue
		}

		index[imageName] = tagName
	}
}

func parseConfig(path string) composeFile {
	file := xos.MustOpen(path)
	defer xos.MustClose(file)
	return parseModel(file)
}

func trimImageName(name string) string {
	if strings.HasPrefix(name, "veupathdb/") {
		return name[10:]
	} else {
		return name
	}
}
