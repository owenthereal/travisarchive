package main

import "fmt"

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

	for _, repo := range repos {
		err := c.DB.Upsert("repos", Query{"lastbuildnumber": repo.LastBuildNumber}, repo)
		if err != nil {
			return err
		}
		fmt.Println(repo.Slug)
	}

	return nil
}
