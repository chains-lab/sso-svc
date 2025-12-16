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

const accountEmailsTable = "account_emails"

type AccountEmail struct {
	AccountID uuid.UUID `db:"account_id"`
	Email     string    `db:"email"`
	Verified  bool      `db:"verified"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type AccountEmailsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewAccountEmails(db *sql.DB) AccountEmailsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return AccountEmailsQ{
		db:       db,
		selector: builder.Select("account_emails.*").From(accountEmailsTable),
		inserter: builder.Insert(accountEmailsTable),
		updater:  builder.Update(accountEmailsTable),
		deleter:  builder.Delete(accountEmailsTable),
		counter:  builder.Select("COUNT(*) AS count").From(accountEmailsTable),
	}
}

func (q AccountEmailsQ) New() AccountEmailsQ {
	return NewAccountEmails(q.db)
}

func (q AccountEmailsQ) Insert(ctx context.Context, input AccountEmail) error {
	values := map[string]interface{}{
		"account_id": input.AccountID,
		"email":      input.Email,
		"verified":   input.Verified,
		"updated_at": input.UpdatedAt,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", accountEmailsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q AccountEmailsQ) Update(ctx context.Context) ([]AccountEmail, error) {
	q.updater = q.updater.
		Set("updated_at", time.Now().UTC()).
		Suffix("RETURNING account_emails.*")

	query, args, err := q.updater.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building update query for %s: %w", accountEmailsTable, err)
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

	var out []AccountEmail
	for rows.Next() {
		var e AccountEmail
		err = rows.Scan(
			&e.AccountID,
			&e.Email,
			&e.Verified,
			&e.UpdatedAt,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning updated account email: %w", err)
		}
		out = append(out, e)
	}

	return out, nil
}

func (q AccountEmailsQ) UpdateEmail(email string) AccountEmailsQ {
	q.updater = q.updater.Set("email", email)
	return q
}

func (q AccountEmailsQ) UpdateVerified(verified bool) AccountEmailsQ {
	q.updater = q.updater.Set("verified", verified)
	return q
}

func (q AccountEmailsQ) Get(ctx context.Context) (AccountEmail, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return AccountEmail{}, fmt.Errorf("building get query for %s: %w", accountEmailsTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var e AccountEmail
	err = row.Scan(
		&e.AccountID,
		&e.Email,
		&e.Verified,
		&e.UpdatedAt,
		&e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AccountEmail{}, nil
		}
		return AccountEmail{}, err
	}

	return e, nil
}

func (q AccountEmailsQ) Select(ctx context.Context) ([]AccountEmail, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", accountEmailsTable, err)
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

	var out []AccountEmail
	for rows.Next() {
		var e AccountEmail
		err = rows.Scan(
			&e.AccountID,
			&e.Email,
			&e.Verified,
			&e.UpdatedAt,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning account_email: %w", err)
		}
		out = append(out, e)
	}

	return out, nil
}

func (q AccountEmailsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", accountEmailsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q AccountEmailsQ) FilterAccountID(accountID uuid.UUID) AccountEmailsQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"account_id": accountID})
	return q
}

func (q AccountEmailsQ) FilterEmail(email string) AccountEmailsQ {
	q.selector = q.selector.Where(sq.Eq{"email": email})
	q.counter = q.counter.Where(sq.Eq{"email": email})
	q.deleter = q.deleter.Where(sq.Eq{"email": email})
	q.updater = q.updater.Where(sq.Eq{"email": email})
	return q
}

func (q AccountEmailsQ) FilterVerified(verified bool) AccountEmailsQ {
	q.selector = q.selector.Where(sq.Eq{"verified": verified})
	q.counter = q.counter.Where(sq.Eq{"verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"verified": verified})
	q.updater = q.updater.Where(sq.Eq{"verified": verified})
	return q
}

func (q AccountEmailsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", accountEmailsTable, err)
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

func (q AccountEmailsQ) Page(limit, offset uint64) AccountEmailsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q AccountEmailsQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
