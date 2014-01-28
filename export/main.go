package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jingweno/travisarchive/db"
)

var (
	execDir  string
	mongoURL string
)

func init() {
	flag.StringVar(&execDir, "e", "", "dir to the mongoexport executable")
	flag.StringVar(&mongoURL, "u", os.Getenv("MONGOHQ_URL"), "URL of the Mongo server")
}

func main() {
	flag.Parse()

	if execDir == "" {
		log.Fatal(fmt.Errorf("specify the dir to the mongoexport executable with -e"))
	}

	if mongoURL == "" {
		log.Fatal(fmt.Errorf("specify the URL of Mongo server with -u"))
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	cols, err := db.DB().CollectionNames()
	if err != nil {
		log.Fatal(err)
	}

	exportBuilds(cols)
}

func exportBuilds(cols []string) {
	oneDayAgo := time.Now().Add(-24 * time.Hour).UTC()
	for _, col := range cols {
		d, err := parseDate(col)
		if err != nil {
			continue
		}

		if d.UTC().After(oneDayAgo) {
			continue
		}

		fmt.Printf("exporting %s...\n", col)
		cmd := &MongoExport{ExecDir: execDir, URL: mongoURL, ColName: col}
		outfile, err := cmd.Run()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("exported %s to %s\n", col, outfile)
	}
}

func parseDate(col string) (time.Time, error) {
	if !strings.HasPrefix(col, "builds_") {
		return time.Time{}, fmt.Errorf("input doesn't include the right prefix")
	}

	timePart := strings.SplitN(col, "_", 2)[1]
	form := "2006_01_02"

	return time.Parse(form, timePart)
}
