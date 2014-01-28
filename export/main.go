package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	col := "builds_2014_01_27"
	cmd := &MongoExport{ExecDir: execDir, URL: mongoURL, ColName: col}
	outfile, err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("exported collection %s to %s\n", col, outfile)
}
