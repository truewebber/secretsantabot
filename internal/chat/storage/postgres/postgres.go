package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/truewebber/secretsantabot/internal/chat/storage"
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

func (p *Postgres) DoOperationOnTx(operation func(context.Context, storage.Tx) error) error {
	ctx := context.Background()

	pgxTx, err := p.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	storageTx := newStorageTx(pgxTx)

	if opErr := operation(ctx, storageTx); opErr != nil {
		if rollbackErr := pgxTx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("rollback tx on do operation: %w: %v", opErr, rollbackErr)
		}

		return fmt.Errorf("do operation: %w", opErr)
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
