package main

import (
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
)

func newServer(
	logger *slog.Logger,
	templates *Template,
	bucket *Bucket,
	subscriber *Subscriber,
	store *Store,
) http.Handler {
	mux := http.NewServeMux()

	public := alice.New(store.sessions.LoadAndSave)
	mux.Handle("GET /", public.Then(newLandingHandler(logger, templates, store)))
	mux.Handle("GET /assets/", public.Then(http.FileServerFS(assets)))
	mux.Handle("GET /products/{link}", public.Then(newDownloadHandler(logger, templates, bucket, store)))
	mux.Handle("POST /subscribe", public.Then(newSubscribeHandler(logger, templates, subscriber)))
	mux.Handle("GET /error", public.Then(newErrorHandler(logger, templates)))
	mux.Handle("GET /login", public.Then(newLoginHandler(logger, templates, store)))
	mux.Handle("POST /login", public.Then(newLoginRequestHandler(logger, templates, store)))
	mux.Handle("GET /logout", public.Then(newLogoutHandler(logger, templates, store)))

	// requires admin
	admin := public.Extend(alice.New(NewAdminOnly(logger, templates, store)))
	mux.Handle("GET /admin", admin.Then(newAdminHandler(logger, templates, bucket, store)))
	mux.Handle("POST /products", admin.Then(newGenerateLinkHandler(logger, templates, store)))

	return mux
}

func NewAdminOnly(l *slog.Logger, t *Template, s *Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.sessions.GetBool(r.Context(), "admin") {
				l.Info("no admin, redirect", "path", r.URL)
				RedirectTo(w, r, "/login")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WriteError(l *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	l.Error(err.Error(), "request", r.URL)
	RedirectTo(w, r, "/error")
}

func RedirectTo(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("Hx-Request") == "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
	} else {
		w.Header().Set("HX-Redirect", path)
		w.Write([]byte{})

	}
}
