package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jingweno/travisarchive/db"
)

var (
	newBuildCrawlerInterval      = 10 * time.Second
	finishedBuildCrawlerInterval = 2 * time.Minute
)

type Crawler interface {
	Crawl()
}

func NewCrawler(travis *Travis, db *db.DB) []Crawler {
	return []Crawler{
		&NewBuildCrawler{travis, db, log.New(os.Stderr, "[NewBuildCrawler] ", log.LstdFlags)},
		&FinishedBuildCrawler{travis, db, log.New(os.Stderr, "[FinishedBuildCrawler] ", log.LstdFlags)},
	}
}

type NewBuildCrawler struct {
	Travis *Travis
	DB     *db.DB
	Logger *log.Logger
}

func (c *NewBuildCrawler) Crawl() {
	ch := time.Tick(newBuildCrawlerInterval)
	for _ = range ch {
		c.Logger.Println("crawling for new builds...")
		c.crawlNewBuilds()
	}
}

func (c *NewBuildCrawler) crawlNewBuilds() {
	repos, err := c.Travis.Repos()
	if err != nil {
		c.Logger.Println(err)
		return
	}

	newBuilds := []string{}
	for _, repo := range repos {
		updated, err := c.DB.Upsert("new_builds", db.Query{"lastbuildid": repo.LastBuildId}, repo)
		if err != nil {
			c.Logger.Println(err)
			continue
		}

		if updated {
			newBuilds = append(newBuilds, repo.Slug)
		}
	}

	c.Logger.Printf("harvested %d builds with %d new builds: %s\n", len(repos), len(newBuilds), strings.Join(newBuilds, ", "))
}

type FinishedBuildCrawler struct {
	Travis *Travis
	DB     *db.DB
	Logger *log.Logger
}

func (c *FinishedBuildCrawler) Crawl() {
	ch := time.Tick(finishedBuildCrawlerInterval)
	for _ = range ch {
		c.Logger.Println("crawling for finsihed builds...")
		c.crawlFinishedBuilds()
	}
}

func (c *FinishedBuildCrawler) crawlFinishedBuilds() {
	colNames, finishedBuilds, skippedBuilds := c.doCrawlFinishedBuilds()
	c.Logger.Printf("fetched %d builds with %d finsihed and %d skipped. Finsihed builds: %s\n", len(finishedBuilds)+len(skippedBuilds), len(finishedBuilds), len(skippedBuilds), strings.Join(finishedBuilds, ", "))

	err := c.ensureColIndexes(colNames)
	if err != nil {
		c.Logger.Println(err)
	}
}

func (c *FinishedBuildCrawler) doCrawlFinishedBuilds() (colNames map[string]string, finishedBuilds []string, skippedBuilds []string) {
	colNames = make(map[string]string)

	var (
		repo *Repo
		//query db.Query
	)
	//query = Query{"lastbuildstartedat": Query{"$gte": oneMinuteAgo()}}
	iter := c.DB.C("new_builds").Find(nil).Sort("lastbuildstartedat").Iter()
	for iter.Next(&repo) {
		build, err := c.crawlFinsihedBuild(repo)
		if err != nil {
			c.Logger.Println(err)
			skippedBuilds = append(skippedBuilds, repo.Slug)
			continue
		}

		colName, updated, err := c.upsertBuild(build)
		if err != nil {
			c.Logger.Println(err)
			continue
		}

		colNames[colName] = colName

		if updated {
			finishedBuilds = append(finishedBuilds, repo.Slug)
		}
	}

	if err := iter.Close(); err != nil {
		c.Logger.Println(err)
	}

	return
}

func (c *FinishedBuildCrawler) crawlFinsihedBuild(repo *Repo) (build *Build, err error) {
	build, err = c.Travis.Build(repo.LastBuildId)
	if err != nil {
		return
	}

	isFinished := build.FinishedAt != nil && build.StartedAt != nil
	if !isFinished {
		err = fmt.Errorf("skipping build: %s - %d\n", repo.Slug, repo.LastBuildId)
		return
	}

	build.Repository = repo

	return
}

func (c *FinishedBuildCrawler) upsertBuild(build *Build) (colName string, updated bool, err error) {
	colName = buildColName(build.StartedAt)

	updated, err = c.DB.Upsert(colName, db.Query{"id": build.Id}, build)
	if err != nil {
		return
	}

	err = c.DB.C("new_builds").Remove(db.Query{"lastbuildid": build.Repository.LastBuildId})
	if err != nil {
		return
	}

	return
}

func (c *FinishedBuildCrawler) ensureColIndexes(colNames map[string]string) error {
	for _, colName := range colNames {
		c.Logger.Printf("ensuring index for collection %s\n", colName)
		err := c.DB.EnsureUniqueIndexKey(colName, "id")
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

func oneMinuteAgo() time.Time {
	now := time.Now()
	return now.Add(-1 * time.Minute).UTC()
}
