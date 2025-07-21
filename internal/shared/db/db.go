package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("Error creating new Postgres pool")
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return pool, nil
}
