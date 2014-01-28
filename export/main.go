package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jingweno/travisarchive/db"
	"github.com/jingweno/travisarchive/export/uploader"
	"github.com/joho/godotenv"
)

var (
	execDir  string
	mongoURL string
)

func init() {
	godotenv.Load("../.env")
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

		log.Printf("exporting %s...\n", col)
		cmd := &MongoExport{ExecDir: execDir, URL: mongoURL, ColName: col}
		outfile, err := cmd.Run()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("exported to %s\n", outfile)

		archiver := &Archiver{outfile}
		outzip, err := archiver.Archive()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("archived to %s\n", outzip)

		zipfile, err := os.Open(outzip)
		if err != nil {
			log.Println(err)
			continue
		}
		defer zipfile.Close()

		filename := fmt.Sprintf("/builds/%s", filepath.Base(outzip))
		u, err := uploader.New("s3")
		if err != nil {
			log.Println(err)
			continue
		}
		err = u.Upload(filename, "application/zip", zipfile)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("uploaded to s3 %s", filename)
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
