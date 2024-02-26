package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/oklog/ulid/v2"
	"github.com/scrot/fourrabbitsstudio/internal/errors"
	"github.com/scrot/fourrabbitsstudio/internal/mail"
	"github.com/scrot/fourrabbitsstudio/internal/storage"
)

func newLandingHandler(l *slog.Logger, t *Template, s *storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := struct {
			User    string
			IsAdmin bool
		}{
			s.Sessions.GetString(r.Context(), "user"),
			false,
		}

		if err := t.RenderPage(w, "landing", data); err != nil {
			l.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
	})
}

func newErrorHandler(l *slog.Logger, t *Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if err := t.RenderPage(w, "error", nil); err != nil {
			l.Error(fmt.Errorf("newErrorHandler: unexpected error: %w", err).Error())
			return
		}
	})
}

func newLoginHandler(l *slog.Logger, t *Template, s *storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := s.Sessions.GetBool(r.Context(), "admin")

		if isAdmin {
			NewRedirect("/admin").ServeHTTP(w, r)
			return
		}
		if err := t.RenderPage(w, "login", nil); err != nil {
			WriteError(l, w, r, err)
			return
		}
	})
}

func newLoginRequestHandler(l *slog.Logger, t *Template, s *storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			WriteError(l, w, r, err)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		// TODO: proper field validation
		if username == "" || password == "" {
			WriteError(l, w, r, errors.ErrMissingField)
			return
		}

		admin, err := s.IsAdmin(r.Context(), username, password)
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		l.Info("login request", "username", username, "admin", admin)

		if admin {
			s.Sessions.Put(r.Context(), "admin", admin)
			s.Sessions.Put(r.Context(), "user", username)
			NewRedirect("/admin").ServeHTTP(w, r)
		} else {
			err := fmt.Errorf("not admin")
			WriteError(l, w, r, err)
			return
		}
	})
}

func newLogoutHandler(l *slog.Logger, t *Template, s *storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := s.Sessions.GetString(r.Context(), "user")

		if username == "" {
			NewRedirect("/").ServeHTTP(w, r)
			return
		}

		s.Sessions.Destroy(r.Context())
		NewRedirect("/").ServeHTTP(w, r)
	})
}

func newSubscribeHandler(l *slog.Logger, t *Template, subscriber *mail.Subscriber) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			WriteError(l, w, r, err)
			return
		}
		email := r.PostFormValue("email")

		if email == "" {
			WriteError(l, w, r, errors.ErrMissingField)
			return
		}

		if err := subscriber.Subscribe(r.Context(), email); err != nil {
			WriteError(l, w, r, err)
			return
		}
		l.Info("new subscriber", "email", email)

		if err := t.RenderPartial(w, "thanks", nil); err != nil {
			WriteError(l, w, r, err)
			return
		}
	})
}

func newSimpleDownloadHandler(l *slog.Logger, s *storage.Store, b storage.Bucket) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		link := r.PathValue("link")
		key, err := s.DownloadLink(r.Context(), link)
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		url, err := b.ObjectURL(r.Context(), key)
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		NewRedirect(url).ServeHTTP(w, r)
	})
}

func newThanksHandler(l *slog.Logger, t *Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := t.RenderPage(w, "thanks", nil); err != nil {
			WriteError(l, w, r, err)
		}
	})
}

func newAdminHandler(l *slog.Logger, t *Template, b storage.Bucket, s *storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("new request", "url", r.URL.String())

		products, err := s.AllProductLinks(r.Context())
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		downloads, err := b.All(r.Context())
		if err != nil {
			WriteError(l, w, r, err)
			return
		}

		for _, download := range downloads {
			var hasProductLink bool

			for _, product := range products {
				if download == product.DownloadLink {
					hasProductLink = true
					break
				}
			}
			if !hasProductLink {
				products = append(products, storage.Product{ProductLink: "", DownloadLink: download})
			}
		}

		for i, p := range products {
			var hasDownloadLink bool
			for _, d := range downloads {
				if p.DownloadLink == d {
					hasDownloadLink = true
					break
				}
			}
			if !hasDownloadLink {
				products[i].DeadLink = true
			}
		}

		l.Info("loaded products", "products", len(products), "downloads", len(downloads))

		data := struct {
			User     string
			IsAdmin  bool
			Products []storage.Product
		}{
			s.Sessions.GetString(r.Context(), "user"),
			s.Sessions.GetBool(r.Context(), "admin"),
			products,
		}

		if err := t.RenderPage(w, "admin", data); err != nil {
			WriteError(l, w, r, err)
			return
		}
	})
}

func newGenerateLinkHandler(l *slog.Logger, t *Template, s *storage.Store) http.Handler {
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

		NewRedirect("/admin").ServeHTTP(w, r)
	})
}
