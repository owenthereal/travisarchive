package main

import "github.com/codegangsta/martini"

func main() {
	m := martini.Classic()
	m.Use(martini.Static("public"))
	m.Run()
}
