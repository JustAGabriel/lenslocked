package views

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
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

func GetTemplate(filename string) (tmpl *template.Template, err error) {
	rootPath, _ := os.Getwd()
	filepath := path.Join(rootPath, viewsSubDirName, templateSubDirName, filename+templatesExtension)
	if exists, err := util.FilepathExists(filepath); !exists {
		errMsg := fmt.Sprintf("template path could not be found/accessed (%s): %v\n", filepath, err)
		err = errors.New(errMsg)
		return tmpl, err
	}
	tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse template (%s): %v\n", filepath, err)
		return nil, errors.New(errMsg)
	}
	return tmpl, err
}

func ParseFS(fs fs.FS, pattern string) (htmlTpl *template.Template, err error) {
	path := path.Join(templateSubDirName, pattern+templatesExtension)
	htmlTpl, err = template.ParseFS(fs, path)
	if err != nil {
		err = fmt.Errorf("parsing template: %w", err)
		return
	}
	return htmlTpl, nil
}

func ExecuteTemplate(w http.ResponseWriter, filename string, data any) {
	log.Println("filename: ", filename)
	tmpl, err := GetTemplate(filename)
	if err != nil {
		log.Printf("could not retrieve template: %v\n", err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	TryWriteTemplate(w, tmpl, data)
}

func TryWriteTemplate(w http.ResponseWriter, tmpl *template.Template, data any) {
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("could not execute template: %v\n", err)
		http.Error(w, "Error while retrieving page.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf8")
}
