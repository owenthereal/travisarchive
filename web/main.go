package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/jingweno/travisarchive/filestore"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func main() {
	fs, err := filestore.New("s3")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	m := martini.Classic()
	m.Map(fs)
	m.Use(martini.Static("../web/public"))
	m.Use(render.Renderer(render.Options{
		Directory: "../web/templates",
		Layout:    "layout",
	}))
	m.Get("/", func(fs *filestore.S3, r render.Render) {
		files, err := fs.List("builds")
		if err != nil {
			log.Fatal(err)
		}

		r.HTML(200, "home", files)
	})

	log.Printf("starting server at %s", port)
	err = http.ListenAndServe(":"+port, m)
	if err != nil {
		log.Fatal(err)
	}
}
