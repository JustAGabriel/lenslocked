package main

import (
	"fmt"
	"net/http"
)

type Router struct{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Home Page. Under construcion...</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Contact Page. Under construcion...</h1>")
}

func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/contact":
		contactHandler(w, r)
	case "/home":
		homeHandler(w, r)
	default:
		http.Error(w, "Page not found.", http.StatusNotFound)
	}
}

func main() {
	http.ListenAndServe(":8000", Router{})
}
