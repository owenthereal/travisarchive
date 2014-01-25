package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Crawler interface {
	Crawl()
}

func NewCrawler(travis *Travis, db *DB) []Crawler {
	return []Crawler{
		&NewBuildCrawler{travis, db, log.New(os.Stderr, "[NewBuildCrawler] ", log.LstdFlags)},
		&FinishedBuildCrawler{travis, db, log.New(os.Stderr, "[FinishedBuildCrawler] ", log.LstdFlags)},
	}
}

type NewBuildCrawler struct {
	Travis *Travis
	DB     *DB
	Logger *log.Logger
}

func (c *NewBuildCrawler) Crawl() {
	ch := time.Tick(10 * time.Second)
	for _ = range ch {
		c.Logger.Printf("crawling for new builds...\n")

		err := c.crawlNewBuilds()
		if err != nil {
			c.Logger.Println(err)
		}
	}
}

func (c *NewBuildCrawler) crawlNewBuilds() (err error) {
	repos, err := c.Travis.Repos()
	if err != nil {
		return err
	}

	newBuilds := []string{}
	for _, repo := range repos {
		updated, err := c.DB.Upsert("new_builds", Query{"lastbuildid": repo.LastBuildId}, repo)
		if err != nil {
			return err
		}

		if updated {
			newBuilds = append(newBuilds, repo.Slug)
		}
	}

	c.Logger.Printf("harvested %d builds with %d new builds: %s\n", len(repos), len(newBuilds), strings.Join(newBuilds, ", "))

	return nil
}

type FinishedBuildCrawler struct {
	Travis *Travis
	DB     *DB
	Logger *log.Logger
}

func (c *FinishedBuildCrawler) Crawl() {
	ch := time.Tick(1 * time.Minute)
	for _ = range ch {
		c.Logger.Printf("crawling for finsihed builds...\n")

		err := c.crawlBuilds()
		if err != nil {
			c.Logger.Println(err)
		}
	}
}

func (c *FinishedBuildCrawler) crawlBuilds() error {
	finishedBuilds := []string{}

	var repo Repo
	colNames := make(map[string]string)
	iter := c.DB.C("new_builds").Find(nil).Iter()
	for iter.Next(&repo) {
		build, err := c.Travis.Build(repo.LastBuildId)
		if err != nil {
			return err
		}

		shouldSkip := build.FinishedAt == nil || build.StartedAt == nil
		var action string
		if shouldSkip {
			action = "skipping"
		} else {
			action = "updating"
		}

		c.Logger.Printf("%s build: %s - %d\n", action, repo.Slug, repo.LastBuildId)

		if shouldSkip {
			continue
		}

		build.Repository = &repo

		colName := buildColName(build.StartedAt)
		colNames[colName] = colName

		updated, err := c.DB.Upsert(colName, Query{"id": build.Id}, build)
		if err != nil {
			return err
		}

		err = c.DB.C("new_builds").Remove(Query{"lastbuildid": repo.LastBuildId})
		if err != nil {
			return err
		}

		if updated {
			finishedBuilds = append(finishedBuilds, repo.Slug)
		}
	}

	c.Logger.Printf("harvested %d finsihed builds: %s\n", len(finishedBuilds), strings.Join(finishedBuilds, ", "))

	return c.ensureColIndexes(colNames)
}

func (c *FinishedBuildCrawler) ensureColIndexes(colNames map[string]string) error {
	for _, colName := range colNames {
		c.Logger.Printf("ensuring index for collection %s\n", colName)
		err := c.DB.EnsureIndexOnField(colName, "id")
		if err != nil {
			return err
		}
	}

	return nil
}

func buildColName(date *time.Time) string {
	buildDate := date.UTC().Format("2006_01_02")
	return fmt.Sprintf("builds_%s", buildDate)
}
