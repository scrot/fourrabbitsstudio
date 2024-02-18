package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

func registerEndpoints(
	mux *http.ServeMux,
	logger *slog.Logger,
	templates *template.Template,
	bucket *Bucket,
) {
	mux.Handle("GET /", newRootHandler(logger, templates))
	mux.Handle("GET /assets/", http.FileServerFS(assets))
	mux.Handle("GET /downloads/", newDownloadHandler(logger, bucket))
}

func newRootHandler(l *slog.Logger, t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request", "url", r.URL.String())
		if err := t.ExecuteTemplate(w, "root.html.tmpl", nil); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})
}

func newDownloadHandler(l *slog.Logger, b *Bucket) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request", "url", r.URL.String())

		objects, err := b.allObjects()
		if err != nil {
			l.Error(fmt.Sprintf("bucket error: %s", err), "bucket", b.name)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(w, strings.Join(objects, "\n"))
	})
}
