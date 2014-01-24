package main

import (
	"log"
	"time"
)

func main() {
	travis := NewTravis("https://api.travis-ci.org")
	crawler := NewCrawler(travis)

	c := time.Tick(10 * time.Second)
	for _ = range c {
		log.Printf("crawling for new repos")

		err := crawler.Crawl()
		if err != nil {
			log.Fatal(err)
		}
	}
}
