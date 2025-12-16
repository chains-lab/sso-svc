package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const accountsTable = "accounts"

type Account struct {
	ID                uuid.UUID `db:"id"`
	Username          string    `db:"username"`
	Role              string    `db:"role"`
	Status            string    `db:"status"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	UsernameUpdatedAt time.Time `db:"username_updated_at"`
}

type AccountsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewAccounts(db *sql.DB) AccountsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return AccountsQ{
		db:       db,
		selector: builder.Select("accounts.*").From(accountsTable),
		inserter: builder.Insert(accountsTable),
		updater:  builder.Update(accountsTable),
		deleter:  builder.Delete(accountsTable),
		counter:  builder.Select("COUNT(*) AS count").From(accountsTable),
	}
}

func (q AccountsQ) New() AccountsQ {
	return NewAccounts(q.db)
}

func (q AccountsQ) Insert(ctx context.Context, input Account) error {
	values := map[string]interface{}{
		"id":                  input.ID,
		"username":            input.Username,
		"role":                input.Role,
		"status":              input.Status,
		"created_at":          input.CreatedAt,
		"updated_at":          input.UpdatedAt,
		"username_updated_at": input.UsernameUpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", accountsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q AccountsQ) Update(ctx context.Context) ([]Account, error) {
	q.updater = q.updater.
		Set("updated_at", time.Now().UTC()).
		Suffix("RETURNING accounts.*")

	query, args, err := q.updater.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building update query for %s: %w", accountsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Account
	for rows.Next() {
		var a Account
		err = rows.Scan(
			&a.ID,
			&a.Username,
			&a.Role,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.UsernameUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning updated account: %w", err)
		}
		out = append(out, a)
	}

	return out, nil
}

func (q AccountsQ) UpdateRole(role string) AccountsQ {
	q.updater = q.updater.Set("role", role)
	return q
}

func (q AccountsQ) UpdateStatus(status string) AccountsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q AccountsQ) UpdateUsername(username string, usernameUpdatedAt time.Time) AccountsQ {
	q.updater = q.updater.
		Set("username", username).
		Set("username_updated_at", usernameUpdatedAt)
	return q
}

func (q AccountsQ) Get(ctx context.Context) (Account, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Account{}, fmt.Errorf("building get query for %s: %w", accountsTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var a Account
	err = row.Scan(
		&a.ID,
		&a.Username,
		&a.Role,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.UsernameUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Account{}, nil
		}
		return Account{}, err
	}

	return a, nil
}

func (q AccountsQ) Select(ctx context.Context) ([]Account, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", accountsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Account
	for rows.Next() {
		var a Account
		err = rows.Scan(
			&a.ID,
			&a.Username,
			&a.Role,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.UsernameUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning account: %w", err)
		}
		out = append(out, a)
	}

	return out, nil
}

func (q AccountsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", accountsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q AccountsQ) FilterID(id uuid.UUID) AccountsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	return q
}

func (q AccountsQ) FilterUsername(username string) AccountsQ {
	q.selector = q.selector.Where(sq.Eq{"username": username})
	q.counter = q.counter.Where(sq.Eq{"username": username})
	q.deleter = q.deleter.Where(sq.Eq{"username": username})
	q.updater = q.updater.Where(sq.Eq{"username": username})
	return q
}

func (q AccountsQ) FilterRole(role string) AccountsQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	return q
}

func (q AccountsQ) FilterStatus(status string) AccountsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	return q
}

func (q AccountsQ) FilterEmail(email string) AccountsQ {
	q.selector = q.selector.
		Join("account_emails ae ON ae.account_id = accounts.id").
		Where(sq.Eq{"ae.email": email})

	q.counter = q.counter.
		Join("account_emails ae ON ae.account_id = accounts.id").
		Where(sq.Eq{"ae.email": email})

	sub := sq.Select("account_id").
		From("account_emails").
		Where(sq.Eq{"email": email})

	q.updater = q.updater.Where(sq.Expr("id IN (?)", sub))
	q.deleter = q.deleter.Where(sq.Expr("id IN (?)", sub))

	return q
}

func (q AccountsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", accountsTable, err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q AccountsQ) Page(limit, offset uint64) AccountsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q AccountsQ) OrderCreatedAt(ascending bool) AccountsQ {
	if ascending {
		q.selector = q.selector.OrderBy("created_at ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC")
	}
	return q
}

func (q AccountsQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
