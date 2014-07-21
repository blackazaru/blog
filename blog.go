package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(
		render.Options{
			Directory: "public",
			Extensions: []string{".tmpl", ".html"},
		},
	))

	m.Get("/", func(w http.ResponseWriter, rr *http.Request,r render.Render) {
			http.Redirect(w, rr, "/home", http.StatusFound)
			r.HTML(200, "home/index", nil)
		})
	m.Get("/about", func(r render.Render) {
			r.HTML(200, "about/index", nil)
		})

	m.Run()
}
