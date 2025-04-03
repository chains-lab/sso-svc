package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/hs-zavet/sso-oauth/internal/service/domain/models"
)

type txKeyType struct{}

var txKey = txKeyType{}

const accountsTable = "accounts"

type Accounts interface {
	New() Accounts

	Insert(ctx context.Context, acc models.Account) error
	Update(ctx context.Context, updates map[string]any) error
	Delete(ctx context.Context) error

	Select(ctx context.Context) ([]models.Account, error)
	Count(ctx context.Context) (int, error)
	Get(ctx context.Context) (*models.Account, error)

	Filter(filters map[string]any) Accounts

	Transaction(fn func(ctx context.Context) error) error

	Page(limit, offset uint64) Accounts
}

type accounts struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewAccounts(db *sql.DB) Accounts {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &accounts{
		db:       db,
		selector: builder.Select("*").From(accountsTable),
		inserter: builder.Insert(accountsTable),
		updater:  builder.Update(accountsTable),
		deleter:  builder.Delete(accountsTable),
		counter:  builder.Select("COUNT(*) AS count").From(accountsTable),
	}
}

func (a *accounts) New() Accounts {
	return NewAccounts(a.db)
}

func (a *accounts) Insert(ctx context.Context, acc models.Account) error {
	values := map[string]interface{}{
		"id":           acc.ID,
		"email":        acc.Email,
		"role":         acc.Role,
		"subscription": acc.Subscription,
		"created_at":   acc.CreatedAt,
		"updated_at":   acc.UpdatedAt,
	}

	query, args, err := a.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for accounts: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = a.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (a *accounts) Update(ctx context.Context, updates map[string]any) error {
	updates["updated_at"] = time.Now().UTC()
	query, args, err := a.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for accounts: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = a.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (a *accounts) Delete(ctx context.Context) error {
	query, args, err := a.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for accounts: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = a.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}

	return nil
}

func (a *accounts) Select(ctx context.Context) ([]models.Account, error) {
	query, args, err := a.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for accounts: %w", err)
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		err := rows.Scan(
			&acc.ID,
			&acc.Email,
			&acc.Role,
			&acc.Subscription,
			&acc.UpdatedAt,
			&acc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning account: %w", err)
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (a *accounts) Count(ctx context.Context) (int, error) {
	query, args, err := a.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for accounts: %w", err)
	}

	var count int
	err = a.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *accounts) Get(ctx context.Context) (*models.Account, error) {
	query, args, err := a.selector.Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("building get query for accounts: %w", err)
	}

	var acc models.Account
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&acc.ID,
		&acc.Email,
		&acc.Role,
		&acc.Subscription,
		&acc.UpdatedAt,
		&acc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (a *accounts) Filter(filters map[string]any) Accounts {
	var validFilters = map[string]bool{
		"id":           true,
		"email":        true,
		"role":         true,
		"subscription": true,
	}
	for key, value := range filters {
		if _, exists := validFilters[key]; !exists {
			continue
		}
		a.selector = a.selector.Where(sq.Eq{key: value})
		a.counter = a.counter.Where(sq.Eq{key: value})
		a.deleter = a.deleter.Where(sq.Eq{key: value})
		a.updater = a.updater.Where(sq.Eq{key: value})
	}
	return a
}

func (a *accounts) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *accounts) Page(limit, offset uint64) Accounts {
	a.counter = a.counter.Limit(limit).Offset(offset)
	a.selector = a.selector.Limit(limit).Offset(offset)
	return a
}
