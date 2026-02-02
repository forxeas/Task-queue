package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	conn *pgxpool.Pool
}

func NewDb(ctx context.Context) (*Db, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}

	return &Db{conn: pool}, nil
}
