package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/repo/pgdb"
)

type Repository struct {
	sql  *pgdb.Queries
	pool *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		sql:  pgdb.New(db),
		pool: db,
	}
}

type txKeyType struct{}

var TxKey = txKeyType{}

func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	return tx, ok
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := TxFromCtx(ctx); ok {
		return fn(ctx)
	}

	tx, err := r.pool.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(ctxWithTx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
