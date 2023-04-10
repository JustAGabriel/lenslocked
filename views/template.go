package views

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/justagabriel/lenslocked/util"
)

const (
	viewsSubDirName    = "views"
	templateSubDirName = "templates"
	templatesExtension = ".html"
)

type Template struct {
	HtmlTmpl template.Template
}

func (t *Template) getTemplate(filename string) (err error) {
	rootPath, _ := os.Getwd()
	filepath := path.Join(rootPath, viewsSubDirName, templateSubDirName, filename+templatesExtension)
	if exists, err := util.FilepathExists(filepath); !exists {
		errMsg := fmt.Sprintf("template path could not be found/accessed (%s): %v\n", filepath, err)
		return errors.New(errMsg)
	}
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse template (%s): %v\n", filepath, err)
		return errors.New(errMsg)
	}
	t.HtmlTmpl = *tmpl
	return nil
}

func ExecuteTemplate(w http.ResponseWriter, filename string, data any) {
	t := Template{}
	log.Println("filename: ", filename)
	err := t.getTemplate(filename)
	if err != nil {
		log.Printf("could not retrieve template: %v\n", err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	err = t.HtmlTmpl.Execute(w, data)
	if err != nil {
		log.Printf("could not execute template: %v\n", err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf8")
}
