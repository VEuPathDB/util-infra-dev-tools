package gh

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Foxcapades/Go-ChainRequest/request/header"
	"github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/sirupsen/logrus"
)

func ListTags(repo string) ([]string, error) {
	logrus.Debugf("making github api request to fetch git tags for repo %s", repo)
	res := simple.GetRequest(repoTagsURL(repo)).
		AddHeader(header.ACCEPT, acceptCT).
		AddHeader(headerApiVersion, apiVersion).
		Submit()
	defer res.Close()

	var tagList []GitHubTag

	if code, err := res.GetResponseCode(); err != nil {
		logrus.Errorf("request to github api failed with error: %s\n", err)
		return nil, err
	} else if code != 200 {
		msg := fmt.Sprintf("got HTTP code %d while requesting tag list for GitHub repo %s", code, repo)
		logrus.Errorln(msg)
		return nil, errors.New(msg)
	}

	if err := res.UnmarshalBody(&tagList, simple.UnmarshallerFunc(json.Unmarshal)); err != nil {
		return nil, err
	}

	out := make([]string, len(tagList))

	for i := range tagList {
		out[i] = tagList[i].Name
	}

	logrus.Debugf("found %d git tags for github repo %s", len(out), repo)

	return out, nil
}

type GitHubTag struct {
	Name   string          `json:"name"`
	Commit GitHubTagCommit `json:"commit"`
	ZipURL string          `json:"zipball_url"`
	TarURL string          `json:"tarball_url"`
	NodeID string          `json:"node_id"`
}

type GitHubTagCommit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}
