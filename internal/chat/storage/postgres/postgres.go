package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	conn *pgxpool.Pool
}

func NewPostgres(connString string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("pgx pool connect: %w", err)
	}

	return &Postgres{conn: pool}, nil
}

func MustNewPostgres(connString string) *Postgres {
	p, err := NewPostgres(connString)
	if err != nil {
		panic(err)
	}

	return p
}
