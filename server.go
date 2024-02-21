package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func newServer(
	logger *slog.Logger,
	templates *Template,
	bucket *Bucket,
	subscriber *Subscriber,
	store *Store,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", newLandingHandler(logger, templates))
	mux.Handle("GET /assets/", http.FileServerFS(assets))
	mux.Handle("GET /products/{link}", newDownloadHandler(logger, templates, bucket, store))
	mux.Handle("POST /subscribe", newSubscribeHandler(logger, templates, subscriber))
	mux.Handle("GET /error", newErrorHandler(logger, templates))

	mux.Handle("GET /login",
		store.sessions.LoadAndSave(newLoginHandler(logger, templates, store)))

	mux.Handle("POST /login",
		store.sessions.LoadAndSave(newLoginRequestHandler(logger, templates, store)))

	// requires admin
	mux.Handle("GET /admin",
		store.sessions.LoadAndSave(adminOnly(newAdminHandler(logger, templates, bucket, store), logger, templates, store)))

	mux.Handle("POST /products",
		store.sessions.LoadAndSave(adminOnly(newGenerateLinkHandler(logger, templates, store), logger, templates, store)))

	return mux
}

func adminOnly(next http.Handler, l *slog.Logger, t *Template, s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.sessions.GetBool(r.Context(), "admin") {
			err := fmt.Errorf("not authorized")
			WriteError(l, w, r, err, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WriteError(l *slog.Logger, w http.ResponseWriter, r *http.Request, err error, msg string) {
	if msg == "" {
		msg = "Whoeps! something went wrong..."
	}

	l.Error(err.Error(), "request", r.URL)
	RedirectTo(w, "/error")
}

func RedirectTo(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Redirect", path)
	w.Write([]byte{})
}
