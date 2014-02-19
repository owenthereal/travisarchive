package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jingweno/travisarchive/db"
	"github.com/jingweno/travisarchive/filestore"
	"github.com/jingweno/travisarchive/util"
	"github.com/joho/godotenv"
)

var (
	execDir  string
	mongoURL string
)

func init() {
	godotenv.Load("../.env")
	flag.StringVar(&execDir, "e", "", "dir to the mongoexport executable")
	flag.StringVar(&mongoURL, "u", os.Getenv("MONGO_URL"), "URL of the Mongo server")
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

	exportBuilds(db, cols)
}

func exportBuilds(db *db.DB, cols []string) {
	oneDayAgo := time.Now().UTC().Add(-24 * time.Hour)
	threeDaysAgo := time.Now().UTC().Add(-3 * 24 * time.Hour)
	for _, col := range cols {
		d, err := util.ParseBuildTime(col)
		if err != nil {
			continue
		}

		if d.UTC().After(oneDayAgo) {
			continue
		}

		outfile, err := exportC(col)
		if err != nil {
			log.Println(err)
			continue
		}

		outzip, err := archiveFile(outfile)
		if err != nil {
			log.Println(err)
			continue
		}

		err = uploadZipfile(outzip)
		if err != nil {
			log.Println(err)
			continue
		}

		// only drop collections that are 3 days old
		if d.UTC().Before(threeDaysAgo) {
			err = dropC(db, col)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func dropC(db *db.DB, col string) error {
	err := db.DropC(col)
	if err != nil {
		return err
	}

	log.Printf("dropped collection %s\n", col)

	return nil
}

func exportC(col string) (string, error) {
	log.Printf("exporting %s...\n", col)
	cmd := &MongoExport{ExecDir: execDir, URL: mongoURL, ColName: col}
	outfile, err := cmd.Run()
	if err != nil {
		return "", err
	}

	log.Printf("exported to %s\n", outfile)

	return outfile, nil
}

func archiveFile(outfile string) (string, error) {
	archiver := &Archiver{outfile}
	outzip, err := archiver.Archive()
	if err != nil {
		return "", err
	}

	log.Printf("archived to %s\n", outzip)

	return outzip, nil
}

func uploadZipfile(outzip string) error {
	zipfile, err := os.Open(outzip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	filename := fmt.Sprintf("/builds/%s", filepath.Base(outzip))
	ds, err := filestore.New("s3")
	if err != nil {
		return err
	}
	err = ds.Upload(filename, "application/zip", zipfile)
	if err != nil {
		return err
	}

	log.Printf("uploaded to s3 %s", filename)

	return nil
}
