package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/google/uuid"
)

const usersPassTable = "users_password"

type UsersPasswordQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewUsersPass(db *sql.DB) UsersPasswordQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return UsersPasswordQ{
		db:       db,
		selector: builder.Select("*").From(usersPassTable),
		inserter: builder.Insert(usersPassTable),
		updater:  builder.Update(usersPassTable),
		deleter:  builder.Delete(usersPassTable),
		counter:  builder.Select("COUNT(*) AS count").From(usersPassTable),
	}
}

func (q UsersPasswordQ) Insert(ctx context.Context, input schemas.UserPasswordModel) error {
	values := map[string]interface{}{
		"user_id":       input.ID,
		"password_hash": input.PassHash,
		"updated_at":    input.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for table %s: %w", usersPassTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UsersPasswordQ) Update(ctx context.Context, input schemas.UserPassUpdateInput) error {
	values := map[string]interface{}{
		"updated_at": input.UpdatedAt,
	}
	if input.PassHash != nil {
		values["password_hash"] = *input.PassHash
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for table %s: %w", usersPassTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UsersPasswordQ) Get(ctx context.Context) (schemas.UserPasswordModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.UserPasswordModel{}, fmt.Errorf("building get query for table %s: %w", usersPassTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var acc schemas.UserPasswordModel
	err = row.Scan(
		&acc.ID,
		&acc.PassHash,
		&acc.UpdatedAt,
	)
	if err != nil {
		return schemas.UserPasswordModel{}, err
	}

	return acc, nil
}

func (q UsersPasswordQ) Select(ctx context.Context) ([]schemas.UserPasswordModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for table %s: %w", usersPassTable, err)
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

	var users []schemas.UserPasswordModel
	for rows.Next() {
		var acc schemas.UserPasswordModel
		err := rows.Scan(
			&acc.ID,
			&acc.PassHash,
			&acc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning table %s: %w", usersPassTable, err)
		}
		users = append(users, acc)
	}

	return users, nil
}

func (q UsersPasswordQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for table %s: %w", usersPassTable, err)
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

func (q UsersPasswordQ) FilterID(userID uuid.UUID) schemas.UsersPasswordQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})

	return q
}

func (q UsersPasswordQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %stable %s: %w", usersPassTable, err)
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

func (q UsersPasswordQ) Page(limit, offset uint) schemas.UsersPasswordQ {
	q.counter = q.counter.Limit(uint64(limit)).Offset(uint64(offset))

	return q
}
