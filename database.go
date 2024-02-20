package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type ProductStore struct {
	conn *pgx.Conn
}

func NewProductStore(ctx context.Context) (*ProductStore, error) {
	dsn, err := Getenv("POSTGRES_DSN")
	if err != nil {
		return nil, err
	}

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
