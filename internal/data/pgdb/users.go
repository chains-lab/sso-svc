package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/google/uuid"
)

const usersTable = "users"

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

func (q UsersQ) Insert(ctx context.Context, input schemas.UserModel) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"role":       input.Role,
		"status":     input.Status,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for users: %w", err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UsersQ) Update(ctx context.Context, input schemas.UserUpdateInput) error {
	values := map[string]any{}

	if input.Status != nil {
		values["status"] = *input.Status
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for users: %w", err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UsersQ) Get(ctx context.Context) (schemas.UserModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.UserModel{}, fmt.Errorf("building get query for users: %w", err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var acc schemas.UserModel
	err = row.Scan(
		&acc.ID,
		&acc.Role,
		&acc.Status,
		&acc.CreatedAt,
	)
	if err != nil {
		return schemas.UserModel{}, err
	}

	return acc, nil
}

func (q UsersQ) Select(ctx context.Context) ([]schemas.UserModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for users: %w", err)
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

	var users []schemas.UserModel
	for rows.Next() {
		var acc schemas.UserModel
		err := rows.Scan(
			&acc.ID,
			&acc.Role,
			&acc.Status,
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
		return fmt.Errorf("building delete query for users: %w", err)
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

func (q UsersQ) FilterID(id uuid.UUID) schemas.UsersQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})

	return q
}

func (q UsersQ) FilterRole(role string) schemas.UsersQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})

	return q
}

func (q UsersQ) FilterStatus(status string) schemas.UsersQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})

	return q
}

func (q UsersQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for users: %w", err)
	}

	var count uint
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

func (q UsersQ) Page(limit, offset uint) schemas.UsersQ {
	q.counter = q.counter.Limit(uint64(limit)).Offset(uint64(offset))

	return q
}
