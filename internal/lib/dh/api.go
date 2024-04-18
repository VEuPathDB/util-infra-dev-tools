package dh

import "fmt"

const (
	namespace = "veupathdb"

	apiURL = "https://hub.docker.com/v2"
)

func makeTagURL(repo, tag string) string {
	return fmt.Sprintf("%s/repositories/%s/%s/tags/%s", apiURL, namespace, repo, tag)
}
