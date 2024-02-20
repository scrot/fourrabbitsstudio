package main

import (
	"embed"
	"html/template"
	"io"
	"path"
)

type Template struct {
	fs       embed.FS
	partials *template.Template
}

func (t *Template) renderPage(w io.Writer, filename string, data any) error {
	templates, err := template.ParseFS(
		t.fs,
		"templates/root.html.tmpl",
		path.Join("templates", filename),
	)
	if err != nil {
		return err
	}

	templates.ExecuteTemplate(w, "root", data)

	return nil
}

func (t *Template) renderPartial(w io.Writer, name string, data any) error {
	return t.partials.ExecuteTemplate(w, name, data)
}
