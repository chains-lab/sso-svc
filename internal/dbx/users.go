package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const usersTable = "users"

type UserModel struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	Verified  bool      `db:"verified,omitempty"`
	Suspended bool      `db:"suspended,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at"`
}

type UserQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewUsers(db *sql.DB) UserQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return UserQ{
		db:       db,
		selector: builder.Select("*").From(usersTable),
		inserter: builder.Insert(usersTable),
		updater:  builder.Update(usersTable),
		deleter:  builder.Delete(usersTable),
		counter:  builder.Select("COUNT(*) AS count").From(usersTable),
	}
}

func (q UserQ) New() UserQ {
	return NewUsers(q.db)
}

func (q UserQ) Insert(ctx context.Context, input UserModel) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"email":      input.Email,
		"role":       input.Role,
		"verified":   input.Verified,
		"suspended":  input.Suspended,
		"updated_at": input.UpdatedAt,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for users: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

type UserUpdateInput struct {
	Role      *string
	Verified  *bool
	Suspended *bool
	UpdatedAt time.Time
}

func (q UserQ) Update(ctx context.Context, input UserUpdateInput) error {
	values := map[string]interface{}{
		"role":       input.Role,
		"verified":   input.Verified,
		"suspended":  input.Suspended,
		"updated_at": input.UpdatedAt,
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for users: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UserQ) Get(ctx context.Context) (UserModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return UserModel{}, fmt.Errorf("building get query for users: %w", err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var acc UserModel
	err = row.Scan(
		&acc.ID,
		&acc.Email,
		&acc.Role,
		&acc.Verified,
		&acc.Suspended,
		&acc.UpdatedAt,
		&acc.CreatedAt,
	)
	if err != nil {
		return UserModel{}, err
	}

	return acc, nil
}

func (q UserQ) Select(ctx context.Context) ([]UserModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for users: %w", err)
	}

	var rows *sql.Rows

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var acc UserModel
		err := rows.Scan(
			&acc.ID,
			&acc.Email,
			&acc.Role,
			&acc.Verified,
			&acc.Suspended,
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

func (q UserQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for users: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}

	return nil
}

func (q UserQ) FilterID(id uuid.UUID) UserQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})

	return q
}

func (q UserQ) FilterEmail(email string) UserQ {
	q.selector = q.selector.Where(sq.Eq{"email": email})
	q.counter = q.counter.Where(sq.Eq{"email": email})
	q.deleter = q.deleter.Where(sq.Eq{"email": email})
	q.updater = q.updater.Where(sq.Eq{"email": email})

	return q
}

func (q UserQ) FilterRole(role string) UserQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})

	return q
}

func (q UserQ) FilterVerified(verified bool) UserQ {
	q.selector = q.selector.Where(sq.Eq{"verified": verified})
	q.counter = q.counter.Where(sq.Eq{"verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"verified": verified})
	q.updater = q.updater.Where(sq.Eq{"verified": verified})

	return q
}

func (q UserQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for users: %w", err)
	}

	var count int64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (q UserQ) Page(limit, offset uint64) UserQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}

func (q UserQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
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
