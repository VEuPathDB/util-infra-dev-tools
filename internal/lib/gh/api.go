package gh

import "fmt"

const (
	apiVersion = "2022-11-28"
	acceptCT   = "application/vnd.github+json"

	headerApiVersion = "X-GitHub-Api-Version"

	apiURL = "https://api.github.com"

	organization = "VEuPathDB"
)

func repoURL(repo string) string {
	return fmt.Sprintf("%s/repos/%s/%s", apiURL, organization, repo)
}

func repoTagsURL(repo string) string {
	return fmt.Sprintf("%s/repos/%s/%s/tags", apiURL, organization, repo)
}

func repoContentURL(repo, filepath string) string {
	return fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiURL, organization, repo, filepath)
}
