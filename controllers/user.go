package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

type Templates struct {
	New template.Template
}

type User struct {
	Templates Templates
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	if err := u.Templates.New.Execute(w, nil); err != nil {
		panic(fmt.Errorf("while parsing 'new user' template: %v", err))
	}
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
	}

	email := r.Form.Get("email")
	pw := r.Form.Get("password")

	fmt.Fprintf(w, "email: %s, pw: %s", email, pw)
}