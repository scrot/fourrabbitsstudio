package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"go.uber.org/automaxprocs/maxprocs"
)

//go:generate npm run build

var (
	port = "8080"

	//go:embed assets
	assets embed.FS

	//go:embed templates
	templateFS embed.FS
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	h := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})
	logger := slog.New(h)

	// Adhere Linux container CPU quota
	procslogger := &procsLogger{logger}
	undo, err := maxprocs.Set(maxprocs.Logger(procslogger.log))
	defer undo()
	if err != nil {
		return fmt.Errorf("failed to set GOMAXPROCS: %w", err)
	}

	partials, err := template.ParseFS(templateFS, "templates/partials/*.tmpl")
	if err != nil {
		return err
	}
	templates := &Template{templateFS, partials}

	bucket, err := NewBucket(ctx)
	if err != nil {
		return err
	}

	subscriber, err := NewSubscriber()
	if err != nil {
		return err
	}

	products, err := NewProductStore(ctx)
	if err != nil {
		return err
	}
	defer products.conn.Close(ctx)

	now, err := products.Now(ctx)
	if err != nil {
		return err
	}
	logger.Info("connected to productstore", "db-time", now.Format(time.RFC822))

	server := newServer(logger, templates, bucket, subscriber, products)

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	logger.Info("server listening", "port", port)
	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

type procsLogger struct {
	Logger *slog.Logger
}

func (l *procsLogger) log(s string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(s, v...))
}
