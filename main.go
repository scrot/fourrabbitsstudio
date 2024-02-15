package main

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

var (
	port = "8080"

	//go:embed assets
	assets embed.FS

	//go:embed *.tmpl
	templates embed.FS
)

func main() {
	http.Handle("GET /assets/", http.FileServerFS(assets))

	t, err := template.ParseFS(templates, "*.tmpl")
	if err != nil {
		slog.Error(err.Error())
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request", "url", r.URL)
		if err := t.ExecuteTemplate(w, "root.html.tmpl", nil); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})

	slog.Info("server started...", "port", port)

	if err := http.ListenAndServe(":"+port, http.DefaultServeMux); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
