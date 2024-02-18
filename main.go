package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

var (
	port = "8080"

	//go:embed assets
	assets embed.FS

	//go:embed templates
	templates embed.FS
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "main: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templates, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		return err
	}

	bucket := NewBucket(ctx, "fourrabbitsstudio")

	server := newServer(logger, templates, bucket)

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	logger.Info("server listening", "port", port)

	return nil
}
