package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"path"
)

const (
	viewsSubDirName    = "views"
	templateSubDirName = "templates"
	templatesExtension = ".html"
)

func ParseFS(fs fs.FS, templates ...string) (htmlTpl *template.Template, err error) {
	var paths []string
	for _, tpl := range templates {
		p := path.Join(templateSubDirName, tpl+templatesExtension)
		paths = append(paths, p)
	}

	htmlTpl, err = template.ParseFS(fs, paths...)
	if err != nil {
		err = fmt.Errorf("parsing template: %w", err)
		return
	}
	return htmlTpl, nil
}
