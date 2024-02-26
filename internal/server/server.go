package server

import (
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"github.com/scrot/fourrabbitsstudio/internal/mail"
	"github.com/scrot/fourrabbitsstudio/internal/storage"
	"github.com/scrot/fourrabbitsstudio/web"
)

func NewServer(
	logger *slog.Logger,
	templates *Template,
	bucket storage.Bucket,
	subscriber *mail.Subscriber,
	store *storage.Store,
) http.Handler {
	mux := http.NewServeMux()

	public := alice.New(store.Sessions.LoadAndSave)
	mux.Handle("GET /assets/", public.Then(http.FileServerFS(web.Assets)))
	mux.Handle("GET /error", public.Then(newErrorHandler(logger, templates)))

	mux.Handle("GET /", public.Then(newLandingHandler(logger, templates, store)))
	mux.Handle("POST /subscribe", public.Then(newSubscribeHandler(logger, templates, subscriber)))
	mux.Handle("GET /products/{link}", public.Then(newSimpleDownloadHandler(logger, store, bucket)))
	mux.Handle("GET /thankyou", public.Then(newThanksHandler(logger, templates)))

	mux.Handle("GET /login", public.Then(newLoginHandler(logger, templates, store)))
	mux.Handle("POST /login", public.Then(newLoginRequestHandler(logger, templates, store)))
	mux.Handle("PUT /login/cancel", public.Then(NewRedirect("/")))
	mux.Handle("GET /logout", public.Then(newLogoutHandler(logger, templates, store)))

	// requires admin
	admin := public.Extend(alice.New(NewAdminOnly(logger, templates, store)))
	mux.Handle("GET /admin", admin.Then(newAdminHandler(logger, templates, bucket, store)))
	mux.Handle("POST /products", admin.Then(newGenerateLinkHandler(logger, templates, store)))

	return mux
}

func NewAdminOnly(l *slog.Logger, t *Template, s *storage.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.Sessions.GetBool(r.Context(), "admin") {
				l.Info("no admin, redirect", "path", r.URL)
				NewRedirect("/login").ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WriteError(l *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	l.Error(err.Error(), "request", r.URL)
	NewRedirect("/error").ServeHTTP(w, r)
}

func NewRedirect(to string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") == "true" {
			w.Header().Set("Hx-Redirect", to)
			w.Write([]byte{})
		} else {
			http.Redirect(w, r, to, http.StatusSeeOther)
		}
	})
}
