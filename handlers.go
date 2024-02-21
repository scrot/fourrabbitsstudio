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

func newLandingHandler(l *slog.Logger, t *Template, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := struct {
			User    string
			IsAdmin bool
		}{
			s.sessions.GetString(r.Context(), "user"),
			false,
		}

		if err := t.renderPage(w, "landing.html.tmpl", data); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})
}

func newErrorHandler(l *slog.Logger, t *Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if err := t.renderPage(w, "error.html.tmpl", nil); err != nil {
			l.Error(fmt.Errorf("newErrorHandler: unexpected error: %w", err).Error())
			return
		}
	})
}

func newLoginHandler(l *slog.Logger, t *Template, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := s.sessions.GetBool(r.Context(), "admin")

		if isAdmin {
			RedirectTo(w, r, "/admin")
			return
		}
		if err := t.renderPage(w, "login.html.tmpl", nil); err != nil {
			WriteError(l, w, r, err)
			return
		}
	})
}

func newLoginRequestHandler(l *slog.Logger, t *Template, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			WriteError(l, w, r, err)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		// TODO: proper field validation
		if username == "" || password == "" {
			WriteError(l, w, r, ErrMissingField)
			return
		}

		admin, err := s.IsAdmin(r.Context(), username, password)
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		l.Info("login request", "username", username, "admin", admin)

		if admin {
			s.sessions.Put(r.Context(), "admin", admin)
			s.sessions.Put(r.Context(), "user", username)
			RedirectTo(w, r, "/admin")
			return
		} else {
			err := fmt.Errorf("not admin")
			WriteError(l, w, r, err)
			return
		}
	})
}

func newLogoutHandler(l *slog.Logger, t *Template, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := s.sessions.GetString(r.Context(), "user")

		if username == "" {
			RedirectTo(w, r, "/")
			return
		}

		s.sessions.Destroy(r.Context())
		RedirectTo(w, r, "/")
	})
}

func newSubscribeHandler(l *slog.Logger, t *Template, subscriber *Subscriber) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			WriteError(l, w, r, err)
			return
		}
		email := r.PostFormValue("email")

		if email == "" {
			WriteError(l, w, r, ErrMissingField)
			return
		}

		if err := subscriber.Subscribe(r.Context(), email); err != nil {
			WriteError(l, w, r, err)
			return
		}
		l.Info("new subscriber", "email", email)

		if err := t.renderPartial(w, "thanks", nil); err != nil {
			WriteError(l, w, r, err)
			return
		}
	})
}

func newDownloadHandler(l *slog.Logger, t *Template, b *Bucket, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("new request", "url", r.URL.String())
		link := r.PathValue("link")

		key, err := s.DownloadLink(r.Context(), link)
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		payload, err := b.downloadObject(r.Context(), key)
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		l.Info("download file", "key", key)

		cd := mime.FormatMediaType("attachment", map[string]string{"filename": key})
		w.Header().Set("Content-Disposition", cd)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeContent(w, r, key, time.Time{}, bytes.NewReader(payload))
	})
}

func newAdminHandler(l *slog.Logger, t *Template, b *Bucket, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("new request", "url", r.URL.String())

		products, err := s.AllProductLinks(r.Context())
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		downloads, err := b.allObjects(r.Context())
		if err != nil {
			WriteError(l, w, r, err)
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

		data := struct {
			User     string
			IsAdmin  bool
			Products []Product
		}{
			s.sessions.GetString(r.Context(), "user"),
			s.sessions.GetBool(r.Context(), "admin"),
			products,
		}

		if err := t.renderPage(w, "admin.html.tmpl", data); err != nil {
			WriteError(l, w, r, err)
			return
		}
	})
}

func newGenerateLinkHandler(l *slog.Logger, t *Template, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			WriteError(l, w, r, err)
			return
		}

		download := r.FormValue("download")
		product := ulid.Make().String()

		l.Info("new link", "download", download, "product", product)
		if err := s.CreateProductLink(r.Context(), product, download); err != nil {
			WriteError(l, w, r, err)
			return
		}

		RedirectTo(w, r, "/admin")
	})
}
