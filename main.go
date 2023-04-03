package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Home Page. Under construcion...</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Contact Page. Under construcion...</h1>")
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/contact" {
		contactHandler(w, r)
		return
	}
	homeHandler(w, r)
}

func main() {
	http.HandleFunc("/", pathHandler)
	http.ListenAndServe(":8000", nil)
}
