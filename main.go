package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := NewDB("mongodb://localhost/travisarchive")
	if err != nil {
		log.Fatal(err)
	}
	err = db.EnsureIndex()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	travis := NewTravis("https://api.travis-ci.org")
	crawlers := NewCrawler(travis, db)

	for _, crawler := range crawlers {
		go crawler.Crawl()
	}

	c := trapSignal()
	<-c
}

func trapSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	return c
}
