package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jingweno/travisarchive/db"
	"github.com/joho/godotenv"
)

var stathat *Stathat

func init() {
	godotenv.Load("../.env")
	stathat = &Stathat{
		StatName: os.Getenv("STATHAT_STAT_NAME"),
		Ezkey:    os.Getenv("STATHAT_EZKEY"),
	}
}

func main() {
	db, err := db.Connect()
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
