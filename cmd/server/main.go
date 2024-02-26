package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/scrot/fourrabbitsstudio/internal/errors"
	"github.com/scrot/fourrabbitsstudio/internal/mail"
	"github.com/scrot/fourrabbitsstudio/internal/server"
	"github.com/scrot/fourrabbitsstudio/internal/storage"
	"github.com/scrot/fourrabbitsstudio/web"
	"go.uber.org/automaxprocs/maxprocs"
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

	templates, error := server.NewTemplate(web.TemplatesFS)
	if error != nil {
		return err
	}

	name := os.Getenv("BUCKET_NAME")
	bucket, err := storage.NewGCBucket(ctx, name)
	if err != nil {
		return err
	}

	subscriber, err := mail.NewSubscriber()
	if err != nil {
		return err
	}

	storeConfig, err := storage.NewStoreConfig(logger)
	if err != nil {
		return err
	}

	store, err := storage.NewStore(ctx, storeConfig)
	if err != nil {
		return err
	}
	defer store.Close()

	now, err := store.Now(ctx)
	if err != nil {
		return err
	}
	logger.Info("connected to productstore", "db-time", now.Format(time.RFC822))

	server := server.NewServer(logger, templates, bucket, subscriber, store)

	host, err := errors.Getenv("HOST")
	if err != nil {
		return err
	}

	port, err := errors.Getenv("PORT")
	if err != nil {
		return err
	}

	httpServer := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}

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
