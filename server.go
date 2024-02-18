package main

import (
	"html/template"
	"log/slog"
	"net/http"
)

func newServer(
	logger *slog.Logger,
	templates *template.Template,
	bucket *Bucket,
) http.Handler {
	mux := http.NewServeMux()
	registerEndpoints(mux, logger, templates, bucket)
	return mux
}
