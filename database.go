package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns          = int32(100)
	defaultMinConns          = int32(0)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

var ErrMissingField = errors.New("missing fields")

type Product struct {
	ProductLink  string
	DownloadLink string
}

type ProductStoreConfig struct {
	*pgxpool.Config
}

func NewProductStoreConfig(l *slog.Logger) (*ProductStoreConfig, error) {
	uname, err := Getenv("POSTGRES_USERNAME")
	if err != nil {
		return nil, err
	}

	secret, err := Getenv("POSTGRES_SECRET")
	if err != nil {
		return nil, err
	}

	host, err := Getenv("POSTGRES_HOST")
	if err != nil {
		return nil, err
	}

	port, err := Getenv("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}

	dbname, err := Getenv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full", uname, secret, host, port, dbname)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = defaultMaxConns
	config.MinConns = defaultMinConns
	config.MaxConnLifetime = defaultMaxConnLifetime
	config.MaxConnIdleTime = defaultMaxConnIdleTime
	config.HealthCheckPeriod = defaultHealthCheckPeriod
	config.ConnConfig.ConnectTimeout = defaultConnectTimeout

	config.BeforeClose = func(_ *pgx.Conn) {
		l.Info("database closed the connection")
	}

	return &ProductStoreConfig{config}, nil
}

type ProductStore struct {
	connPool *pgxpool.Pool
}

func NewProductStore(ctx context.Context, config *ProductStoreConfig) (*ProductStore, error) {
	connPool, err := pgxpool.NewWithConfig(ctx, config.Config)
	if err != nil {
		return nil, fmt.Errorf("NewProductStore: load configuration: %w", err)
	}

	conn, err := connPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewProductStore: aquire connection: %w", err)
	}
	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("NewProductStore: pinging database: %w", err)
	}

	return &ProductStore{connPool}, nil
}

func (s *ProductStore) Close() {
	s.connPool.Close()
}

func (s *ProductStore) Now(ctx context.Context) (time.Time, error) {
	const stmt = `SELECT NOW()`

	var now time.Time
	if err := s.connPool.QueryRow(ctx, stmt).Scan(&now); err != nil {
		return time.Time{}, nil
	}

	return now, nil
}

func (s *ProductStore) All(ctx context.Context) ([]Product, error) {
	const stmt = `
  SELECT product_link, download_link
  FROM products
  `
	rows, err := s.connPool.Query(ctx, stmt)
	if err != nil {
		return []Product{}, err
	}

	var ps []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ProductLink, &p.DownloadLink); err != nil {
			return []Product{}, err
		}
		ps = append(ps, p)
	}

	return ps, nil
}

func (s *ProductStore) CreateLink(ctx context.Context, productLink, downloadLink string) error {
	const stmt = `
  INSERT INTO products (product_link, download_link)
  VALUES ($1, $2)
  `

	if productLink == "" || downloadLink == "" {
		return ErrMissingField
	}

	if _, err := s.connPool.Exec(ctx, stmt, productLink, downloadLink); err != nil {
		return err
	}

	return nil
}

func (s *ProductStore) DownloadLink(ctx context.Context, productLink string) (string, error) {
	const stmt = `
  SELECT download_link
  FROM products
  WHERE product_link = $1
  `

	if productLink == "" {
		return "", ErrMissingField
	}

	var downloadLink string
	if err := s.connPool.QueryRow(ctx, stmt, productLink).Scan(&downloadLink); err != nil {
		return "", err
	}

	return downloadLink, nil
}
