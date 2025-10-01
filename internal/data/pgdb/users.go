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

const usersTable = "users"

type User struct {
	ID     uuid.UUID `db:"id"`
	Role   string    `db:"role"`
	Status string    `db:"status"`

	PasswordHash string    `db:"password_hash"`
	PasswordUpAt time.Time `db:"password_updated_at"`

	Email    string `db:"email"`
	EmailVer bool   `db:"email_verified"`

	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type UsersQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewUsers(db *sql.DB) UsersQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return UsersQ{
		db:       db,
		selector: builder.Select("*").From(usersTable),
		inserter: builder.Insert(usersTable),
		updater:  builder.Update(usersTable),
		deleter:  builder.Delete(usersTable),
		counter:  builder.Select("COUNT(*) AS count").From(usersTable),
	}
}

func (q UsersQ) New() UsersQ {
	return NewUsers(q.db)
}

func (q UsersQ) Insert(ctx context.Context, input User) error {
	values := map[string]interface{}{
		"id":     input.ID,
		"role":   input.Role,
		"status": input.Status,

		"password_hash":       input.PasswordHash,
		"password_updated_at": input.PasswordUpAt,

		"email":          input.Email,
		"email_verified": input.EmailVer,

		"updated_at": input.UpdatedAt,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", usersTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UsersQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", usersTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UsersQ) UpdateStatus(status string) UsersQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q UsersQ) UpdatePassword(passwordHash string, passwordUpAt time.Time) UsersQ {
	q.updater = q.updater.
		Set("password_hash", passwordHash).
		Set("password_updated_at", passwordUpAt)
	return q
}

func (q UsersQ) UpdateEmailVerified(emailVer bool) UsersQ {
	q.updater = q.updater.Set("email_verified", emailVer)
	return q
}

func (q UsersQ) UpdateEmail(email string) UsersQ {
	q.updater = q.updater.Set("email", email)
	return q
}

func (q UsersQ) Get(ctx context.Context) (User, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return User{}, fmt.Errorf("building get query for %s: %w", usersTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var acc User
	err = row.Scan(
		&acc.ID,
		&acc.Role,
		&acc.Status,
		&acc.PasswordHash,
		&acc.PasswordUpAt,
		&acc.Email,
		&acc.EmailVer,
		&acc.UpdatedAt,
		&acc.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}

		return User{}, err
	}

	return acc, nil
}

func (q UsersQ) Select(ctx context.Context) ([]User, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", usersTable, err)
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

	var users []User
	for rows.Next() {
		var acc User
		err = rows.Scan(
			&acc.ID,
			&acc.Role,
			&acc.Status,
			&acc.PasswordHash,
			&acc.PasswordUpAt,
			&acc.Email,
			&acc.EmailVer,
			&acc.UpdatedAt,
			&acc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}
		users = append(users, acc)
	}

	return users, nil
}

func (q UsersQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", usersTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}

	return nil
}

func (q UsersQ) FilterID(id uuid.UUID) UsersQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})

	return q
}

func (q UsersQ) FilterEmail(email string) UsersQ {
	q.selector = q.selector.Where(sq.Eq{"email": email})
	q.counter = q.counter.Where(sq.Eq{"email": email})
	q.deleter = q.deleter.Where(sq.Eq{"email": email})
	q.updater = q.updater.Where(sq.Eq{"email": email})

	return q
}

func (q UsersQ) FilterRole(role string) UsersQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})

	return q
}

func (q UsersQ) FilterStatus(status string) UsersQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})

	return q
}

func (q UsersQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", usersTable, err)
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

func (q UsersQ) Page(limit, offset uint64) UsersQ {
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q UsersQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
