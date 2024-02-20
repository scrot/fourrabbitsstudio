package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func newLandingHandler(l *slog.Logger, t *Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("new request", "url", r.URL.String())
		if err := t.renderPage(w, "landing.html.tmpl", nil); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})
}

func newSubscribeHandler(l *slog.Logger, t *Template, subscriber *Subscriber) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}

		email := r.PostFormValue("email")
		l.With("email", email)

		if err := subscriber.Subscribe(r.Context(), email); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
		l.Info("new subscriber")

		if err := t.renderPartial(w, "thanks", nil); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}
	})
}

func newDownloadHandler(l *slog.Logger, b *Bucket, ps *ProductStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("new request", "url", r.URL.String())
	})
}

func newAdminHandler(l *slog.Logger, t *Template, b *Bucket, ps *ProductStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("new request", "url", r.URL.String())

		products, err := ps.All(r.Context())
		if err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}

		objects, err := b.allObjects(r.Context())
		if err != nil {
			l.Error(fmt.Sprintf("bucket error: %s", err), "bucket", b.name)
			return
		}

		linked := make(map[string]struct{})
		for _, pos := range products {
			for _, po := range pos {
				linked[po] = struct{}{}
			}
		}

		var unlinked []string
		for _, o := range objects {
			if _, ok := linked[o]; !ok {
				unlinked = append(unlinked, o)
			}
		}

		data := struct {
			Products map[string][]string
			Unlinked []string
		}{products, unlinked}

		l.Info("loaded products", "products", len(products), "objects", len(objects), "unlinked", len(unlinked))

		if err := t.renderPage(w, "admin.html.tmpl", data); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}
	})
}

func newGenerateLinkHandler(l *slog.Logger, t *Template, ps *ProductStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		object := r.PathValue("object")
		l.Info("generate new link", "object", object)

		link := uuid.New()
		if err := ps.Link(r.Context(), link.String(), object); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/admin")
		w.Write([]byte{})
	})
}
