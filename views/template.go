package views

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/justagabriel/lenslocked/util"
)

const (
	templateSubDirName = "templates"
	templatesExtension = ".html"
)

func ExecuteTemplate(w http.ResponseWriter, filename string, data any) {
	filepath := path.Join(templateSubDirName, filename, templatesExtension)
	if exists, err := util.FilepathExists(filepath); !exists {
		log.Printf("template path could not be found/accessed (%s): %v\n", filepath, err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("could not parse template (%s): %v\n", filepath, err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("could not execute template (%s): %v\n", filepath, err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf8")
}
