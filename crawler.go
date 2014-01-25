package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func NewCrawler(travis *Travis, db *DB) []Crawler {
	return []Crawler{
		&NewBuildCrawler{travis, db},
		&BuildCrawler{travis, db},
	}
}

type Crawler interface {
	Crawl()
}

type NewBuildCrawler struct {
	Travis *Travis
	DB     *DB
}

func (c *NewBuildCrawler) Crawl() {
	ch := time.Tick(10 * time.Second)
	for _ = range ch {
		log.Printf("crawling for new builds...")

		err := c.crawlNewBuilds()
		if err != nil {
			log.Fatal(err)
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

	log.Printf("harvested %d builds with %d new builds: %s\n", len(repos), len(newBuilds), strings.Join(newBuilds, ", "))

	return nil
}

type BuildCrawler struct {
	Travis *Travis
	DB     *DB
}

func (c *BuildCrawler) Crawl() {
	ch := time.Tick(10 * time.Second)
	for _ = range ch {
		log.Printf("crawling for builds...")

		err := c.crawlBuilds()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *BuildCrawler) crawlBuilds() error {
	finishedBuilds := []string{}

	var repo Repo
	colNames := make(map[string]string)
	iter := c.DB.C("new_builds").Find(nil).Iter()
	for iter.Next(&repo) {
		build, err := c.Travis.Build(repo.LastBuildId)
		if err != nil {
			return err
		}

		var (
			status     string
			shouldSkip bool
		)
		if build.FinishedAt == nil || build.StartedAt == nil {
			shouldSkip = true
			status = "skipping"
		} else {
			status = "updating"
		}

		log.Printf("fetched build for %s - %d...%s", repo.Slug, repo.LastBuildId, status)

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

	log.Printf("inserted %d builds: %s\n", len(finishedBuilds), strings.Join(finishedBuilds, ", "))

	for _, colName := range colNames {
		log.Printf("ensuring index for collection %s\n", colName)
		err := c.DB.EnsureIndexOnField(colName, "id")
		if err != nil {
			return err
		}
	}

	return nil
}

func buildColName(date *time.Time) string {
	buildDate := date.Format("2006-01-02")
	buildDate = strings.Replace(buildDate, "-", "_", -1)
	return fmt.Sprintf("builds_%s", buildDate)
}
