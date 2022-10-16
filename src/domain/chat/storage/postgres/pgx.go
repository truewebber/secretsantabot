package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/truewebber/secretsantabot/domain/chat/storage"
)

type pgxAdapter struct {
	conn *pgxpool.Pool
}

func NewPGX(connString string) (storage.Storage, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("pgxAdapter pool connect: %w", err)
	}

	return &pgxAdapter{conn: pool}, nil
}

func MustNewPGX(connString string) storage.Storage {
	p, err := NewPGX(connString)
	if err != nil {
		panic(err)
	}

	return p
}

func (p *pgxAdapter) DoOperationOnTx(
	ctx context.Context,
	operation func(context.Context, storage.Tx) error,
) error {
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

func (p *pgxAdapter) DoLockedOperationOnTx(
	ctx context.Context,
	lockID int64,
	operation func(context.Context, storage.Tx) error,
) error {
	doErr := p.DoOperationOnTx(ctx, func(opCtx context.Context, tx storage.Tx) error {
		if err := tx.LockTx(ctx, lockID); err != nil {
			return fmt.Errorf("lock on tx: %w", err)
		}

		if opErr := operation(ctx, tx); opErr != nil {
			return fmt.Errorf("do operation after lock: %w", opErr)
		}

		return nil
	})

	if doErr != nil {
		return fmt.Errorf("do operation with lock on tx: %w", doErr)
	}

	return nil
}
