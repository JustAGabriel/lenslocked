package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justagabriel/lenslocked/controllers"
	"github.com/justagabriel/lenslocked/util"
	"github.com/justagabriel/lenslocked/views"
)

const (
	baseLayoutFilename = "tailwind"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tpl := util.Must(views.ParseFS(views.FS, baseLayoutFilename, "home"))
	r.Get("/home", controllers.StaticHandler(tpl))

	tpl = util.Must(views.ParseFS(views.FS, baseLayoutFilename, "contact"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = util.Must(views.ParseFS(views.FS, baseLayoutFilename, "faq"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe(":8000", r)
}
