package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

const usersTable = "users"

type UserModel struct {
	ID           uuid.UUID  `db:"id"`
	Email        string     `db:"email"`
	Role         roles.Role `db:"role"`
	Subscription uuid.UUID  `db:"subscription"`
	UpdatedAt    *time.Time `db:"updated_at,omitempty"`
	CreatedAt    time.Time  `db:"created_at"`
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

func (a UserQ) New() UserQ {
	return NewUsers(a.db)
}

type UserInsertInput struct {
	ID           uuid.UUID
	Email        string
	Role         roles.Role
	Subscription uuid.UUID
	CreatedAt    time.Time
}

func (a UserQ) Insert(ctx context.Context, input UserInsertInput) error {
	values := map[string]interface{}{
		"id":           input.ID,
		"email":        input.Email,
		"role":         input.Role,
		"subscription": input.Subscription,
		"created_at":   input.CreatedAt,
	}

	query, args, err := a.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for users: %w", err)
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

type UserUpdateInput struct {
	Role         *roles.Role
	Subscription *uuid.UUID
	UpdatedAt    time.Time
}

func (a UserQ) Update(ctx context.Context, input UserUpdateInput) error {
	values := map[string]interface{}{
		"role":         input.Role,
		"subscription": input.Subscription,
		"updated_at":   input.UpdatedAt,
	}

	query, args, err := a.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for users: %w", err)
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

func (a UserQ) Delete(ctx context.Context) error {
	query, args, err := a.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for users: %w", err)
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

func (a UserQ) Select(ctx context.Context) ([]UserModel, error) {
	query, args, err := a.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for users: %w", err)
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
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
			&acc.Subscription,
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

func (a UserQ) Count(ctx context.Context) (int, error) {
	query, args, err := a.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for users: %w", err)
	}

	var count int
	err = a.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a UserQ) Get(ctx context.Context) (UserModel, error) {
	query, args, err := a.selector.Limit(1).ToSql()
	if err != nil {
		return UserModel{}, fmt.Errorf("building get query for users: %w", err)
	}

	var acc UserModel
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&acc.ID,
		&acc.Email,
		&acc.Role,
		&acc.Subscription,
		&acc.UpdatedAt,
		&acc.CreatedAt,
	)
	if err != nil {
		return UserModel{}, err
	}

	return acc, nil
}

func (a UserQ) FilterID(id uuid.UUID) UserQ {
	a.selector = a.selector.Where(sq.Eq{"id": id})
	a.counter = a.counter.Where(sq.Eq{"id": id})
	a.deleter = a.deleter.Where(sq.Eq{"id": id})
	a.updater = a.updater.Where(sq.Eq{"id": id})
	return a
}

func (a UserQ) FilterEmail(email string) UserQ {
	a.selector = a.selector.Where(sq.Eq{"email": email})
	a.counter = a.counter.Where(sq.Eq{"email": email})
	a.deleter = a.deleter.Where(sq.Eq{"email": email})
	a.updater = a.updater.Where(sq.Eq{"email": email})
	return a
}

func (a UserQ) FilterRole(role roles.Role) UserQ {
	a.selector = a.selector.Where(sq.Eq{"role": role})
	a.counter = a.counter.Where(sq.Eq{"role": role})
	a.deleter = a.deleter.Where(sq.Eq{"role": role})
	a.updater = a.updater.Where(sq.Eq{"role": role})
	return a
}

func (a UserQ) FilterSubscription(subscription uuid.UUID) UserQ {
	a.selector = a.selector.Where(sq.Eq{"subscription": subscription})
	a.counter = a.counter.Where(sq.Eq{"subscription": subscription})
	a.deleter = a.deleter.Where(sq.Eq{"subscription": subscription})
	a.updater = a.updater.Where(sq.Eq{"subscription": subscription})
	return a
}

func (a UserQ) Transaction(fn func(ctx context.Context) error) error {
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

func (a UserQ) Page(limit, offset uint64) UserQ {
	a.counter = a.counter.Limit(limit).Offset(offset)
	a.selector = a.selector.Limit(limit).Offset(offset)
	return a
}

func (a UserQ) Drop(ctx context.Context) error {
	query, args, err := a.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building drop query for users: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = a.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("error executing drop query: %w", err)
	}

	return nil
}
