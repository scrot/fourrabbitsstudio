package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrMissingField = errors.New("missing fields")

type Product struct {
	ProductLink  string
	DownloadLink string
}

type ProductStore struct {
	conn *pgx.Conn
}

func NewProductStore(ctx context.Context) (*ProductStore, error) {
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
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &ProductStore{conn}, nil
}

func (s *ProductStore) Now(ctx context.Context) (time.Time, error) {
	const stmt = `SELECT NOW()`

	var now time.Time
	if err := s.conn.QueryRow(ctx, stmt).Scan(&now); err != nil {
		return time.Time{}, nil
	}

	return now, nil
}

func (s *ProductStore) All(ctx context.Context) ([]Product, error) {
	const stmt = `
  SELECT product_link, download_link
  FROM products
  `
	rows, err := s.conn.Query(ctx, stmt)
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

	if _, err := s.conn.Exec(ctx, stmt, productLink, downloadLink); err != nil {
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
	if err := s.conn.QueryRow(ctx, stmt, productLink).Scan(&downloadLink); err != nil {
		return "", err
	}

	return downloadLink, nil
}
