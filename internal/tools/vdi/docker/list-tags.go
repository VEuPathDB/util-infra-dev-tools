package vdi_docker

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"vpdb-dev-tool/internal/lib/col"
	"vpdb-dev-tool/internal/lib/dh"
	"vpdb-dev-tool/internal/lib/gh"
	"vpdb-dev-tool/internal/lib/xio"
	"vpdb-dev-tool/internal/tools/vdi/projects"
)

func mainGenTagger(opts cliOpts) {
	logrus.Debugf("generating tagger yaml entries for the latest available tagged VDI docker images")

	if opts.writeToVersionsFile {
		tagsToFile()
	} else {
		tagsToStdout()
	}
}

func tagsToFile() {
	stat := xio.MustStat("versions.yml")

	original := xio.MustOpen("versions.yml", os.O_RDONLY, stat.Mode().Perm())
	defer xio.QuietCloseFile(original)

	images, err := taggerList()
	if err != nil {
		logrus.Fatalln(err)
	}

	temp := xio.ReqRWFile{File: xio.MustOpen("versions.yml~", os.O_WRONLY|os.O_CREATE|os.O_EXCL, stat.Mode().Perm())}
	defer xio.QuietCloseFile(temp.File)

	scan := bufio.NewScanner(original)

	for scan.Scan() {
		line := scan.Text()
		idx := strings.IndexByte(line, ':')

		if idx < 1 {
			_, _ = temp.WriteLine(line)
			continue
		}

		key := line[:idx]

		if !images.Contains(key) {
			_, _ = temp.WriteLine(line)
			continue
		}

		_, _ = temp.WriteString(key)
		_, _ = temp.WriteString(": ")
		_, _ = temp.WriteLine(images.Require(key))

		images.Delete(key)
	}

	if scan.Err() != nil {
		logrus.Fatalf("encountered error while scanning versions.yml: %s", scan.Err())
	}

	// Catch any dangling images in case the versions file didn't have them all
	images.ForEach(func(img string, tag string) {
		_, _ = temp.WriteString(img)
		_, _ = temp.WriteString(": ")
		_, _ = temp.WriteLine(tag)
	})

	xio.QuietCloseFile(temp.File)
	if err := os.Rename("versions.yml~", "versions.yml"); err != nil {
		logrus.Fatalf("failed renaming versions.yml~ to versions.yml: %s", err)
	}
}

func tagsToStdout() {
	images, err := taggerList()
	if err != nil {
		logrus.Fatalln(err)
	}

	images.ForEach(func(img string, tag string) { fmt.Printf("%s: %s\n", img, tag) })
}

func taggerList() (col.OrderedMap[string, string], error) {
	projects := vdi_projects.DeployableImageProducers()
	out := col.NewOrderedMap[string, string](len(projects))

	for _, project := range projects {
		tags, err := gh.ListTags(project.Name)

		if err != nil {
			return nil, fmt.Errorf("encountered error while attempting to fetch tags for project %s: %s", project, err)
		}

		for _, image := range project.Images {
			imageTag := "latest"

			for _, tag := range tags {
				exists := false

				if strings.HasPrefix(tag, "v") {
					tTag := tag[1:]

					exists, err = dh.TestTag(image, tTag)

					if err != nil {
						return nil, fmt.Errorf("encountered error while testing dockerhub for image %s:%s: %s", image, tTag, err)
					} else if exists {
						imageTag = tTag
					}
				}

				if !exists {
					exists, err = dh.TestTag(image, tag)

					if err != nil {
						return nil, fmt.Errorf("encountered error while testing dockerhub for image %s:%s: %s", image, tag, err)
					} else if exists {
						imageTag = tag
					}
				}

				if exists {
					break
				}
			}

			out.Put(image, imageTag)
		}
	}

	return out, nil
}
