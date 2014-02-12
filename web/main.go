package main

import (
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/jingweno/travisarchive/filestore"
	"github.com/joho/godotenv"
)

type ByTime []filestore.File

func (t ByTime) Len() int           { return len(t) }
func (t ByTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTime) Less(i, j int) bool { return t[i].Time.After(t[j].Time) }

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
		sort.Sort(ByTime(files))
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
