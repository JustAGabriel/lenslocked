package controllers

import (
	"html/template"
	"net/http"
)

type Static struct {
	Template template.Template
}

func StaticHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}
