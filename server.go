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
	subscriber *Subscriber,
	products *ProductStore,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", newRootHandler(logger, templates))
	mux.Handle("GET /assets/", http.FileServerFS(assets))
	mux.Handle("GET /downloads/", newDownloadHandler(logger, bucket, products))
	mux.Handle("POST /subscribe", newSubscribeHandler(logger, templates, subscriber))
	return mux
}
