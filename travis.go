package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func NewTravis(host string) *Travis {
	return &Travis{Host: host, Client: &http.Client{}}
}

type Travis struct {
	Host   string
	Client *http.Client
}

func (t *Travis) Repos() (repos []Repo, err error) {
	url := fmt.Sprintf("%s/repos", t.Host)
	err = t.get(url, &repos)

	return
}

func (t *Travis) Build(id int) (build *Build, err error) {
	url := fmt.Sprintf("%s/builds/%d", t.Host, id)
	err = t.get(url, &build)
	return
}

func (t *Travis) get(url string, output interface{}) error {
	return t.do("GET", url, nil, output)
}

func (t *Travis) do(verb, url string, input interface{}, output interface{}) (err error) {
	var reader io.Reader
	if input != nil {
		b, e := json.Marshal(input)
		if e != nil {
			err = e
			return
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(verb, url, reader)
	if err != nil {
		return
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%s %s - %d: %s", verb, url, resp.StatusCode, string(body))
		return
	}

	err = json.Unmarshal(body, output)

	return
}

type Repo struct {
	Id                  int        `json:"id"`
	Slug                string     `json:"slug"`
	Description         string     `json:"description"`
	LastBuildFinishedAt *time.Time `json:"last_build_finished_at"`
	LastBuildStartedAt  *time.Time `json:"last_build_started_at"`
	LastBuildLanguage   string     `json:"last_build_language"`
	LastBuildId         int        `json:"last_build_id"`
	LastBuildNumber     string     `json:"last_build_number"`
	LastBuildStatus     int        `json:"last_build_status"`
	LastBuildResult     int        `json:"last_build_result"`
	LastBuildDuration   int        `json:"last_build_result"`
}

type Build struct {
	Id             int                      `json:"id"`
	RepositoryId   int                      `json:"repository_id"`
	Number         string                   `json:"number"`
	Config         map[string]interface{}   `json:"config"`
	State          string                   `json:"state"`
	Result         int                      `json:"result"`
	Status         int                      `json:"status"`
	StartedAt      *time.Time               `json:"started_at"`
	FinishedAt     *time.Time               `json:"finished_at"`
	Duration       int                      `json:"duration"`
	Commit         string                   `json:"commit"`
	Branch         string                   `json:"branch"`
	Message        string                   `json:"message"`
	CommittedAt    *time.Time               `json:"committed_at"`
	AuthorName     string                   `json:"author_name"`
	AuthorEmail    string                   `json:"author_email"`
	CommitterName  string                   `json:"committer_name"`
	CommitterEmail string                   `json:"committer_email"`
	CompareUrl     string                   `json:"compare_url"`
	EventType      string                   `json:"event_type"`
	Matrix         []map[string]interface{} `json:"matrix"`
	Repository     *Repo
}
