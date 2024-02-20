package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

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
