package main

import "fmt"

func NewCrawler(travis *Travis) *Crawler {
	return &Crawler{travis}
}

type Crawler struct {
	Travis *Travis
}

func (c *Crawler) Crawl() error {
	repos, err := c.Travis.Repos()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Println(repo.Slug)
	}

	return nil
}
