package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

func newRootHandler(l *slog.Logger, t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request", "url", r.URL.String())
		if err := t.ExecuteTemplate(w, "root.html.tmpl", nil); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})
}

func newSubscribeHandler(l *slog.Logger, t *template.Template, subscriber *Subscriber) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}

		email := r.PostFormValue("email")
		l.With("email", email)

		if err := subscriber.Subscribe(r.Context(), email); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
		l.Info("new subscriber")

		if err := t.ExecuteTemplate(w, "thanks", nil); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})
}

func newDownloadHandler(l *slog.Logger, b *Bucket, ps *ProductStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request", "url", r.URL.String())

		objects, err := b.allObjects(r.Context())
		if err != nil {
			l.Error(fmt.Sprintf("bucket error: %s", err), "bucket", b.name)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(w, strings.Join(objects, "\n"))
	})
}
