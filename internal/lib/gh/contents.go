package gh

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Foxcapades/Go-ChainRequest/request/header"
	"github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/sirupsen/logrus"
)

func GetPrivateFileContents(repo, filepath, branch, token string) ([]byte, error) {
	logrus.Debugf("making an authenticated github api request to fetch file %s from repo %s", filepath, repo)

	res := simple.GetRequest(repoContentURL(repo, filepath)+"?ref="+branch).
		AddHeader(header.AUTHORIZATION, "Bearer "+token).
		AddHeader(header.ACCEPT, acceptCT).
		AddHeader(headerApiVersion, apiVersion).
		Submit()
	defer res.Close()

	if code, err := res.GetResponseCode(); err != nil {
		logrus.Errorf("request to github failed with error: %s", err)
		return nil, err
	} else if code != 200 {
		msg := fmt.Sprintf("got HTTP code %d while requesting file %s from GitHub repo %s", code, filepath, repo)
		logrus.Errorln(msg)
		return nil, errors.New(msg)
	}

	var content FileContent

	if err := res.UnmarshalBody(&content, simple.UnmarshallerFunc(json.Unmarshal)); err != nil {
		return nil, err
	}

	logrus.Debugf("successfully fetched file %s from GitHub repo %s", filepath, repo)

	return content.Content, nil
}

type FileContent struct {
	Type        string       `json:"type"`
	Encoding    string       `json:"encoding"`
	Size        uint64       `json:"size"`
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	Content     []byte       `json:"content"`
	SHA         string       `json:"sha"`
	URL         string       `json:"url"`
	GitURL      string       `json:"git_url"`
	HTMLURL     string       `json:"html_url"`
	DownloadURL string       `json:"download_url"`
	Links       ContentLinks `json:"_links"`
}

type ContentLinks struct {
	Git  string `json:"git"`
	Self string `json:"self"`
	HTML string `json:"html"`
}
