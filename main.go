package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justagabriel/lenslocked/views"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "dude"
	}
	views.ExecuteTemplate(w, "home", struct{ Name string }{Name: name})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	views.ExecuteTemplate(w, "contact", nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	views.ExecuteTemplate(w, "faq", nil)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/contact", contactHandler)
	r.Get("/home", homeHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe(":8000", r)
}
