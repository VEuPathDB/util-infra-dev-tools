package vdi_docker

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"vpdb-dev-tool/internal/lib/dh"
	"vpdb-dev-tool/internal/lib/gh"
	"vpdb-dev-tool/internal/tools/vdi/projects"
)

func mainGenTagger() {
	logrus.Debugf("generating tagger yaml entries for the latest available tagged VDI docker images")
	images, err := taggerList()
	if err != nil {
		logrus.Fatalln(err)
	}

	for image, tag := range images {
		fmt.Printf("%s: %s\n", image, tag)
	}
}

func taggerList() (map[string]string, error) {
	projects := vdi_projects.DeployableImageProducers()
	out := make(map[string]string, len(projects))

	for project, images := range projects {
		tags, err := gh.ListTags(project)

		if err != nil {
			return nil, fmt.Errorf("encountered error while attempting to fetch tags for project %s: %s", project, err)
		}

		for _, image := range images {
			imageTag := "latest"

			for _, tag := range tags {
				exists := false

				if strings.HasPrefix(tag, "v") {
					tTag := tag[1:]

					exists, err = dh.TestTag(image, tTag)

					if err != nil {
						return nil, fmt.Errorf("encountered error while testing dockerhub for image %s:%s: %s", image, tTag, err)
					}
				}

				if !exists {
					exists, err = dh.TestTag(image, tag)

					if err != nil {
						return nil, fmt.Errorf("encountered error while testing dockerhub for image %s:%s: %s", image, tag, err)
					}
				}

				if exists {
					imageTag = tag
					break
				}
			}

			out[image] = imageTag
		}
	}

	return out, nil
}
