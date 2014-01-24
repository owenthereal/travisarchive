package main

import (
	"log"
	"time"
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
	crawler := NewCrawler(travis, db)

	c := time.Tick(10 * time.Second)
	for _ = range c {
		log.Printf("crawling for new repos")

		err := crawler.Crawl()
		if err != nil {
			log.Fatal(err)
		}
	}
}
