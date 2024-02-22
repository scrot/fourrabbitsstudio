package main

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"
)

type Template struct {
	pages    map[string]func(io.Writer, any) error
	partials map[string]func(io.Writer, any) error
}

func NewTemplate(templateFS embed.FS) (*Template, error) {
	tmpl := &Template{
		pages:    make(map[string]func(io.Writer, any) error),
		partials: make(map[string]func(io.Writer, any) error),
	}

	pagePaths, err := fs.Glob(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, pagePath := range pagePaths {
		t, err := template.ParseFS(templateFS, "templates/root.html.tmpl", pagePath)
		if err != nil {
			return nil, err
		}

		k := key(pagePath)
		tmpl.pages[k] = rexecute(t)
	}

	partialPaths, err := fs.Glob(templateFS, "templates/partials/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, partialPath := range partialPaths {
		t, err := template.ParseFS(templateFS, partialPath)
		if err != nil {
			return nil, err
		}

		k := key(partialPath)
		tmpl.partials[k] = t.Execute
	}

	return tmpl, nil
}

var ErrNotFound = errors.New("template not found")

func (t *Template) RenderPage(w io.Writer, name string, data any) error {
	render, ok := t.pages[name]
	if !ok {
		return ErrNotFound
	}

	var buf bytes.Buffer
	err := render(&buf, data)
	if err != nil {
		return err
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (t *Template) RenderPartial(w io.Writer, name string, data any) error {
	render, ok := t.partials[name]
	if !ok {
		return ErrNotFound
	}

	var buf bytes.Buffer
	err := render(&buf, data)
	if err != nil {
		return err
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func rexecute(t *template.Template) func(io.Writer, any) error {
	return func(w io.Writer, data any) error {
		return t.ExecuteTemplate(w, "root", data)
	}
}

func key(p string) string {
	_, file := path.Split(p)
	name, _ := strings.CutSuffix(file, ".html.tmpl")
	return name
}
