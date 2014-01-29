package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/martini"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func main() {
	port := os.Getenv("PORT")
	m := martini.Classic()
	m.Use(martini.Static("public"))

	log.Printf("starting server at %s", port)
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}
