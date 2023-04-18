package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justagabriel/lenslocked/controllers"
	"github.com/justagabriel/lenslocked/util"
	"github.com/justagabriel/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tpl := util.Must(views.GetTemplate("home"))
	r.Get("/home", controllers.StaticHandler(tpl))

	tpl = util.Must(views.GetTemplate("contact"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = util.Must(views.GetTemplate("faq"))
	r.Get("/faq", controllers.StaticHandler(tpl))

	tpl = util.Must(views.GetTemplate("legal"))
	r.Get("/legal", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe(":8000", r)
}
