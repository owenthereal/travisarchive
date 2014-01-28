package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type MongoExport struct {
	ExecDir string
	URL     string
	ColName string
}

func (e *MongoExport) Run() (outfile string, err error) {
	bin, err := e.findBin()
	if err != nil {
		return
	}

	dir, err := ioutil.TempDir("", "mongoexport")
	if err != nil {
		return
	}

	uu, err := url.Parse(e.URL)
	if err != nil {
		return
	}

	outfile = filepath.Join(dir, e.ColName+".json")
	args := e.buildArgs(uu, outfile)

	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func (e *MongoExport) findBin() (bin string, err error) {
	execDir, err := filepath.Abs(e.ExecDir)
	if err != nil {
		return
	}

	bin = fmt.Sprintf("mongoexport_%s_%s", runtime.GOOS, runtime.GOARCH)
	bin = filepath.Join(execDir, bin)

	return
}

func (e *MongoExport) buildArgs(uu *url.URL, outfile string) (results []string) {
	host := uu.Host
	if strings.Contains(host, ":") {
		ss := strings.Split(host, ":")
		host = ss[0]
		port := ss[1]
		results = append(results, "--port", port)
	}
	results = append(results, "-h", host)

	user := uu.User
	if user != nil {
		pass, _ := user.Password()
		results = append(results, "-u", uu.User.Username())
		results = append(results, "-p", pass)
	}

	db := strings.TrimPrefix(uu.Path, "/")
	results = append(results, "-d", db)

	results = append(results, "-c", e.ColName)
	results = append(results, "--out", outfile)

	return
}
