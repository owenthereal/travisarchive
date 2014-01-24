package main

import (
	"log"
	"strings"
)

func NewCrawler(travis *Travis, db *DB) *Crawler {
	return &Crawler{travis, db}
}

type Crawler struct {
	Travis *Travis
	DB     *DB
}

func (c *Crawler) Crawl() error {
	repos, err := c.Travis.Repos()
	if err != nil {
		return err
	}

	newBuilds := []string{}
	for _, repo := range repos {
		updated, err := c.DB.Upsert("new_builds", Query{"lastbuildid": repo.LastBuildID}, repo)
		if err != nil {
			return err
		}

		if updated {
			newBuilds = append(newBuilds, repo.Slug)
		}
	}

	log.Printf("harvested %d builds with %d new builds: %s\n", len(repos), len(newBuilds), strings.Join(newBuilds, ", "))

	return nil
}
