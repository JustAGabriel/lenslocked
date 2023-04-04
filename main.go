package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "dude"
	}
	fmt.Fprintf(w, "<h1>Home Page. Under construcion...Sorry, %s</h1>", name)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Contact Page. Under construcion...</h1>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	pageHTML := `<h1>FAQ Page. Under Construction...</h1>`
	fmt.Fprint(w, pageHTML)
}

func main() {
	r := chi.NewRouter()
	r.Get("/contact", contactHandler)
	r.Get("/home", homeHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe(":8000", r)
}
