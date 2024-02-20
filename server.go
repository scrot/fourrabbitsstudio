package main

import (
	"log/slog"
	"net/http"
)

func newServer(
	logger *slog.Logger,
	templates *Template,
	bucket *Bucket,
	subscriber *Subscriber,
	products *ProductStore,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", newLandingHandler(logger, templates))
	mux.Handle("GET /assets/", http.FileServerFS(assets))
	mux.Handle("GET /downloads/{link}", newDownloadHandler(logger, bucket, products))
	mux.Handle("POST /subscribe", newSubscribeHandler(logger, templates, subscriber))

	// requires admin
	mux.Handle("GET /admin", newAdminHandler(logger, templates, bucket, products))
	mux.Handle("PUT /objects/{object}", newGenerateLinkHandler(logger, templates, products))
	return mux
}
