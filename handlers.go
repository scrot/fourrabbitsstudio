package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

func newLandingHandler(l *slog.Logger, t *Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

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

		if err := subscriber.Subscribe(r.Context(), email); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
		l.Info("new subscriber", "email", email)

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
		link := r.PathValue("link")

		key, err := ps.DownloadLink(r.Context(), link)
		if err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}

		payload, err := b.downloadObject(r.Context(), key)
		if err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}

		l.Info("download file", "key", key)

		cd := mime.FormatMediaType("attachment", map[string]string{"filename": key})
		w.Header().Set("Content-Disposition", cd)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeContent(w, r, key, time.Time{}, bytes.NewReader(payload))
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

		downloads, err := b.allObjects(r.Context())
		if err != nil {
			l.Error(fmt.Sprintf("bucket error: %s", err), "bucket", b.name)
			return
		}

		for _, download := range downloads {
			var seen bool
			for _, product := range products {
				if download == product.DownloadLink {
					seen = true
					break
				}
			}
			if !seen {
				products = append(products, Product{ProductLink: "", DownloadLink: download})
			}
		}

		l.Info("loaded products", "products", len(products), "downloads", len(downloads))

		if err := t.renderPage(w, "admin.html.tmpl", products); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}
	})
}

func newGenerateLinkHandler(l *slog.Logger, t *Template, ps *ProductStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}

		download := r.FormValue("download")
		product := ulid.Make().String()

		l.Info("new link", "download", download, "product", product)
		if err := ps.CreateLink(r.Context(), product, download); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/admin")
		w.Write([]byte{})
	})
}
