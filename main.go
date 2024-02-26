package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/crypto/bcrypt"
)

//go:generate tailwindcss build -i ./assets/base.css -o assets/tailwind.css

var (
	port = "8080"

	//go:embed assets
	assets embed.FS

	//go:embed templates
	templatesFS embed.FS
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

	templates, error := NewTemplate(templatesFS)
	if error != nil {
		return err
	}

	bucket, err := NewBucket(ctx)
	if err != nil {
		return err
	}

	subscriber, err := NewSubscriber()
	if err != nil {
		return err
	}

	storeConfig, err := NewStoreConfig(logger)
	if err != nil {
		return err
	}

	store, err := NewStore(ctx, storeConfig)
	if err != nil {
		return err
	}
	defer store.Close()

	now, err := store.Now(ctx)
	if err != nil {
		return err
	}
	logger.Info("connected to productstore", "db-time", now.Format(time.RFC822))

	server := newServer(logger, templates, bucket, subscriber, store)

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	secret, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	logger.Info("server listening", "port", port, "secret", string(secret))
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
