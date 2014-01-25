package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
)

type TestServer struct {
	*httptest.Server
	*http.ServeMux
}

func setupTravisServer() *TestServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &TestServer{Server: server, ServeMux: mux}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func respondWithJSON(w http.ResponseWriter, s string) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	fmt.Fprint(w, s)
}

func loadFixture(f string) string {
	pwd, _ := os.Getwd()
	p := filepath.Join(pwd, "fixtures", f)
	c, _ := ioutil.ReadFile(p)
	return string(c)
}

func TestTravis_Repos(t *testing.T) {
	server := setupTravisServer()
	defer server.Close()

	server.HandleFunc("/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("repos.json"))
	})

	travis := NewTravis(server.URL)
	repos, _ := travis.Repos()

	assert.Equal(t, 25, len(repos))

	repo := repos[0]
	assert.Equal(t, 1584783, repo.Id)
	assert.Equal(t, "Vayleryn/VaylerynLib", repo.Slug)
	assert.Equal(t, 17567797, repo.LastBuildId)
}

func TestTravis_Build(t *testing.T) {
	server := setupTravisServer()
	defer server.Close()

	server.HandleFunc("/builds/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("build.json"))
	})

	travis := NewTravis(server.URL)
	build, _ := travis.Build(1)

	assert.Equal(t, 17567141, build.Id)
	assert.Equal(t, 1584783, build.RepositoryId)
	assert.Equal(t, 61, build.Duration)
}
